// load measures acknowledged WebSocket edits against one contended board.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syncframe/internal/model"
	"time"
)

type Result struct {
	Users           int      `json:"users"`
	Attempted       int      `json:"attempted"`
	Succeeded       int      `json:"succeeded"`
	Failed          int      `json:"failed"`
	Seconds         float64  `json:"seconds"`
	EventsPerSecond float64  `json:"eventsPerSecond"`
	P50MS           float64  `json:"p50Ms"`
	P95MS           float64  `json:"p95Ms"`
	P99MS           float64  `json:"p99Ms"`
	Converged       bool     `json:"converged"`
	Errors          []string `json:"errors"`
}

func main() {
	base := flag.String("url", "http://127.0.0.1:8080", "server URL")
	steps := flag.String("users", "1,5,20", "comma-separated users")
	events := flag.Int("events", 30, "updates per user (1-50)")
	output := flag.String("output", "artifacts/load-results.json", "JSON report")
	flag.Parse()
	if *events < 1 || *events > 50 {
		panic("events must be 1-50 for the per-connection edit rate policy")
	}
	report := struct {
		At             string   `json:"at"`
		Environment    string   `json:"environment"`
		Transport      string   `json:"transport"`
		UpdatesPerUser int      `json:"updatesPerUser"`
		Results        []Result `json:"results"`
	}{At: time.Now().UTC().Format(time.RFC3339), Environment: fmt.Sprintf("%s %s/%s; CPUs=%d; target=%s", runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), *base), Transport: "WebSocket; one board per stage; create + sequential acknowledged updates; no think time; includes SQL commit; not a capacity ceiling", UpdatesPerUser: *events}
	failed := false
	for _, raw := range strings.Split(*steps, ",") {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 100 {
			panic("users must be 1-100")
		}
		r := run(*base, n, *events)
		report.Results = append(report.Results, r)
		b, _ := json.Marshal(r)
		fmt.Println(string(b))
		if r.Failed > 0 || !r.Converged {
			failed = true
		}
	}
	b, _ := json.MarshalIndent(report, "", "  ")
	if err := os.MkdirAll("artifacts", 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(*output, b, 0644); err != nil {
		panic(err)
	}
	if failed {
		os.Exit(1)
	}
}
func run(base string, users, updates int) Result {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Timeout: 10 * time.Second}
	resp, err := client.Post(base+"/api/boards", "application/json", bytes.NewBufferString(`{"title":"Measured WebSocket workload"}`))
	if err != nil {
		panic(err)
	}
	var created struct {
		Board model.State `json:"board"`
	}
	err = json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()
	if err != nil || resp.StatusCode != 201 {
		panic("create load board failed")
	}
	id := created.Board.ID
	req, _ := http.NewRequest("GET", base+"/api/boards/"+id, nil)
	headers := http.Header{"Origin": []string{base}}
	for _, c := range jar.Cookies(req.URL) {
		headers.Add("Cookie", c.String())
	}
	conns := make([]*websocket.Conn, users)
	for i := range conns {
		url := "ws" + strings.TrimPrefix(base, "http") + "/api/boards/" + id + fmt.Sprintf("/ws?name=Load-%d&after=0", i)
		conn, _, err := websocket.DefaultDialer.Dial(url, headers)
		if err != nil {
			panic(err)
		}
		conns[i] = conn
		defer conn.Close()
	}
	result := Result{Users: users, Attempted: users * (updates + 1), Errors: []string{}}
	var mu sync.Mutex
	var wg sync.WaitGroup
	latencies := []float64{}
	start := time.Now()
	for i, conn := range conns {
		wg.Add(1)
		go func(i int, c *websocket.Conn) {
			defer wg.Done()
			for j := 0; j <= updates; j++ {
				objectID := fmt.Sprintf("user-object-%d", i)
				op := model.Operation{ID: model.ID(), Actor: fmt.Sprintf("Load-%d", i), ObjectID: objectID}
				if j == 0 {
					op.Kind = "create"
					op.Object = &model.Object{ID: objectID, Kind: "rect", Width: 100, Height: 80, Color: "#dce7ff"}
				} else {
					x := float64(j)
					op.Kind = "update"
					op.Patch = &model.Patch{X: &x}
				}
				t0 := time.Now()
				_ = c.SetWriteDeadline(t0.Add(10 * time.Second))
				err := c.WriteJSON(map[string]any{"type": "operation", "operation": op})
				_ = c.SetReadDeadline(t0.Add(10 * time.Second))
				if err == nil {
					for {
						var msg struct {
							Type        string `json:"type"`
							OperationID string `json:"operationId"`
							Error       string `json:"error"`
						}
						if err = c.ReadJSON(&msg); err != nil {
							break
						}
						if msg.OperationID == op.ID {
							if msg.Type == "rejected" {
								err = fmt.Errorf("rejected: %s", msg.Error)
							}
							break
						}
					}
				}
				mu.Lock()
				if err != nil {
					result.Failed += updates + 1 - j
					result.Errors = append(result.Errors, fmt.Sprintf("user %d edit %d: %v", i, j, err))
					mu.Unlock()
					return
				}
				result.Succeeded++
				latencies = append(latencies, float64(time.Since(t0).Microseconds())/1000)
				mu.Unlock()
			}
		}(i, conn)
	}
	wg.Wait()
	result.Seconds = time.Since(start).Seconds()
	result.EventsPerSecond = float64(result.Succeeded) / result.Seconds
	sort.Float64s(latencies)
	percentile := func(p float64) float64 {
		if len(latencies) == 0 {
			return 0
		}
		return latencies[int(float64(len(latencies)-1)*p)]
	}
	result.P50MS = percentile(.5)
	result.P95MS = percentile(.95)
	result.P99MS = percentile(.99)
	response, err := client.Get(base + "/api/boards/" + id)
	if err == nil {
		raw, _ := io.ReadAll(response.Body)
		response.Body.Close()
		var state model.State
		if json.Unmarshal(raw, &state) == nil {
			result.Converged = state.Revision == int64(result.Attempted) && len(state.Objects) == users
			for i := 0; i < users; i++ {
				if state.Objects[fmt.Sprintf("user-object-%d", i)].X != float64(updates) {
					result.Converged = false
				}
			}
		}
	}
	return result
}
