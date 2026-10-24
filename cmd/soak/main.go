// soak is a bounded verification client, not an application service.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"syncframe/internal/model"
)

type peer struct {
	conn    *websocket.Conn
	mu      sync.Mutex
	state   model.State
	pending map[string]chan string
	err     string
}

func (p *peer) read() {
	for {
		var m struct {
			Type        string      `json:"type"`
			State       model.State `json:"state"`
			Event       model.Event `json:"event"`
			OperationID string      `json:"operationId"`
			Error       string      `json:"error"`
		}
		if err := p.conn.ReadJSON(&m); err != nil {
			p.mu.Lock()
			p.err = err.Error()
			p.mu.Unlock()
			return
		}
		p.mu.Lock()
		switch m.Type {
		case "snapshot":
			p.state = m.State
		case "event":
			if err := model.Replay(&p.state, m.Event); err != nil {
				p.err = err.Error()
			}
		case "ack", "rejected":
			if ch := p.pending[m.OperationID]; ch != nil {
				ch <- m.Error
				delete(p.pending, m.OperationID)
			}
		}
		p.mu.Unlock()
	}
}
func (p *peer) send(op model.Operation) error {
	ch := make(chan string, 1)
	p.mu.Lock()
	p.pending[op.ID] = ch
	p.mu.Unlock()
	p.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err := p.conn.WriteJSON(map[string]any{"type": "operation", "operation": op}); err != nil {
		return err
	}
	select {
	case e := <-ch:
		if e != "" {
			return fmt.Errorf("%s", e)
		}
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("ACK timeout")
	}
}
func main() {
	base := flag.String("url", "http://127.0.0.1:8080", "real server")
	users := flag.Int("users", 20, "1-100 independent sockets on one board")
	duration := flag.Duration("duration", 5*time.Minute, "1s-5m active edit window, excluding setup")
	interval := flag.Duration("interval", time.Second, "per-user update interval; missed ticks are not queued")
	updates := flag.Int("updates", 0, "1-50 for an unpaced finite burst; 0 selects sustained mode")
	output := flag.String("output", "artifacts/portfolio/soak.json", "report path")
	flag.Parse()
	if *users < 1 || *users > 100 || *duration < time.Second || *duration > 5*time.Minute || *updates < 0 || *updates > 50 || *interval < 250*time.Millisecond || int64(*users)*(int64(*duration / *interval)+2) > 90000 {
		panic("bounded workload limits exceeded")
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Timeout: 15 * time.Second}
	r, err := client.Post(*base+"/api/boards", "application/json", bytes.NewBufferString(`{"title":"Bounded sustained verification"}`))
	must(err)
	var created struct {
		Board model.State `json:"board"`
	}
	must(json.NewDecoder(r.Body).Decode(&created))
	r.Body.Close()
	if r.StatusCode != 201 {
		panic("create failed")
	}
	id := created.Board.ID
	req, _ := http.NewRequest("GET", *base+"/api/boards/"+id, nil)
	headers := http.Header{"Origin": []string{*base}}
	for _, c := range jar.Cookies(req.URL) {
		headers.Add("Cookie", c.String())
	}
	peers := make([]*peer, *users)
	for i := range peers {
		c, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(*base, "http")+"/api/boards/"+id+fmt.Sprintf("/ws?name=Soak-%d&after=0", i), headers)
		must(err)
		peers[i] = &peer{conn: c, pending: map[string]chan string{}}
		go peers[i].read()
		defer c.Close()
	}
	// State-driven readiness, bounded by a deadline.
	readyUntil := time.Now().Add(15 * time.Second)
	for {
		ready := true
		for _, p := range peers {
			p.mu.Lock()
			ready = ready && p.state.ID == id
			p.mu.Unlock()
		}
		if ready {
			break
		}
		if time.Now().After(readyUntil) {
			panic("initial snapshots missing")
		}
		time.Sleep(20 * time.Millisecond)
	}
	start := time.Now()
	end := start.Add(*duration)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var latencies []float64
	attempted, failed := 0, 0
	errors := []string{}
	last := make([]int, *users)
	for i, p := range peers {
		wg.Add(1)
		go func(i int, p *peer) {
			defer wg.Done()
			tick := time.NewTicker(*interval)
			defer tick.Stop()
			for j := 0; ; j++ {
				if *updates > 0 && j > *updates {
					return
				}
				if j > 0 && *updates == 0 {
					select {
					case <-tick.C:
					case <-time.After(time.Until(end)):
						return
					}
					if !time.Now().Before(end) {
						return
					}
				}
				oid := fmt.Sprintf("soak-object-%d", i)
				op := model.Operation{ID: model.ID(), Actor: fmt.Sprintf("Soak-%d", i), ObjectID: oid}
				if j == 0 {
					op.Kind = "create"
					op.Object = &model.Object{ID: oid, Kind: "rect", Width: 100, Height: 80, Color: "#dce7ff"}
				} else {
					x := float64(j)
					op.Kind = "update"
					op.Patch = &model.Patch{X: &x}
				}
				t0 := time.Now()
				err := p.send(op)
				mu.Lock()
				attempted++
				if err != nil {
					failed++
					errors = append(errors, fmt.Sprintf("user %d edit %d: %v", i, j, err))
					mu.Unlock()
					return
				}
				latencies = append(latencies, float64(time.Since(t0).Microseconds())/1000)
				last[i] = j
				mu.Unlock()
			}
		}(i, p)
	}
	wg.Wait()
	seconds := time.Since(start).Seconds()
	var state model.State
	get := func(path string, out any) {
		r, e := client.Get(*base + "/api/boards/" + id + path)
		must(e)
		defer r.Body.Close()
		if r.StatusCode != 200 {
			panic("state/events request failed")
		}
		must(json.NewDecoder(r.Body).Decode(out))
	}
	get("", &state)
	replay := model.State{ID: id, Title: state.Title, Objects: map[string]model.Object{}}
	seen := map[string]bool{}
	unique := true
	for {
		var es []model.Event
		get(fmt.Sprintf("/events?after=%d", replay.Revision), &es)
		for _, e := range es {
			if seen[e.Operation.ID] {
				unique = false
			}
			seen[e.Operation.ID] = true
			must(model.Replay(&replay, e))
		}
		if len(es) < 500 {
			break
		}
	}
	expected := state.Revision == int64(attempted-failed) && len(state.Objects) == *users
	if *updates > 0 {
		expected = expected && attempted == *users*(*updates+1)
	}
	for i, x := range last {
		o, ok := state.Objects[fmt.Sprintf("soak-object-%d", i)]
		expected = expected && ok && o.X == float64(x)
	}
	until := time.Now().Add(30 * time.Second)
	converged := false
	for {
		converged = true
		for _, p := range peers {
			p.mu.Lock()
			converged = converged && reflect.DeepEqual(p.state, state) && p.err == ""
			p.mu.Unlock()
		}
		if converged || time.Now().After(until) {
			break
		}
		time.Sleep(25 * time.Millisecond)
	}
	sort.Float64s(latencies)
	q := func(p float64) float64 {
		if len(latencies) == 0 {
			return 0
		}
		return latencies[int(float64(len(latencies)-1)*p)]
	}
	report := map[string]any{"at": start.UTC().Format(time.RFC3339), "environment": fmt.Sprintf("%s %s/%s; CPUs=%d; GOMAXPROCS=%d", runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.GOMAXPROCS(0)), "target": *base, "boardId": id, "users": *users, "requestedSeconds": duration.Seconds(), "intervalMs": interval.Milliseconds(), "seconds": seconds, "attempted": attempted, "succeeded": len(latencies), "failed": failed, "failureRate": float64(failed) / float64(attempted), "eventsPerSecond": float64(len(latencies)) / seconds, "p50Ms": q(.5), "p95Ms": q(.95), "p99Ms": q(.99), "expectedState": expected, "allWebSocketsConverged": converged, "replayMatchesState": reflect.DeepEqual(replay, state), "uniqueOperationIds": unique, "finalRevision": state.Revision, "errors": errors, "method": "create then one outstanding ACK per user, paced updates; continuous socket readers; no reconnect/retry; setup excluded; actual attempts counted; not rendering latency"}
	if *updates > 0 {
		report["updatesPerUser"] = *updates
		report["requestedSeconds"] = nil
		report["intervalMs"] = 0
		report["method"] = "unpaced create + finite sequential updates per user; one outstanding ACK; continuous socket readers; no reconnect/retry; setup excluded; actual attempts counted; not rendering latency"
	}
	must(os.MkdirAll(filepath.Dir(*output), 0755))
	b, _ := json.MarshalIndent(report, "", "  ")
	must(os.WriteFile(*output, b, 0644))
	fmt.Println(string(b))
	if failed > 0 || !expected || !converged || !unique || !reflect.DeepEqual(replay, state) {
		os.Exit(1)
	}
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
