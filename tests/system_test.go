package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"syncframe/internal/model"
	"syncframe/internal/pipeline"
	"syncframe/internal/store"
)

func config(t *testing.T) (string, *http.Client) {
	t.Helper()
	base := os.Getenv("TEST_BASE_URL")
	if base == "" {
		t.Skip("TEST_BASE_URL not set: real HTTP system tests skipped")
	}
	jar, _ := cookiejar.New(nil)
	return base, &http.Client{Jar: jar, Timeout: 8 * time.Second}
}
func request(t *testing.T, c *http.Client, method, path string, body any, status int, out any) {
	t.Helper()
	raw, _ := json.Marshal(body)
	req, err := http.NewRequest(method, path, bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != status {
		t.Fatalf("%s %s: %d wanted %d: %s", method, path, resp.StatusCode, status, b)
	}
	if out != nil {
		if err = json.Unmarshal(b, out); err != nil {
			t.Fatal(err)
		}
	}
}
func create(t *testing.T, c *http.Client, base string) (model.State, string) {
	var result struct {
		Board model.State `json:"board"`
		Token string      `json:"token"`
	}
	request(t, c, "POST", base+"/api/boards", map[string]string{"title": "System verification"}, 201, &result)
	return result.Board, result.Token
}
func operation(i int) model.Operation {
	id := fmt.Sprintf("shape-%d", i)
	return model.Operation{ID: model.ID(), Actor: "System test", Kind: "create", ObjectID: id, Object: &model.Object{ID: id, Kind: "rect", Width: 150, Height: 100, Color: "#dce7ff"}}
}
func edit(t *testing.T, c *http.Client, base, id string, op model.Operation) model.Event {
	var result struct {
		Event model.Event `json:"event"`
	}
	request(t, c, "POST", base+"/api/boards/"+id+"/operations", op, 200, &result)
	return result.Event
}
func ws(t *testing.T, c *http.Client, base, id string, after int64) *websocket.Conn {
	t.Helper()
	u, _ := http.NewRequest("GET", base+"/api/boards/"+id, nil)
	h := http.Header{}
	for _, cookie := range c.Jar.Cookies(u.URL) {
		h.Add("Cookie", cookie.String())
	}
	h.Set("Origin", base)
	url := "ws" + strings.TrimPrefix(base, "http") + "/api/boards/" + id + fmt.Sprintf("/ws?name=System%%20test&after=%d", after)
	conn, _, err := websocket.DefaultDialer.Dial(url, h)
	if err != nil {
		t.Fatal(err)
	}
	return conn
}
func waitRevision(t *testing.T, c *websocket.Conn, want int64) {
	t.Helper()
	c.SetReadDeadline(time.Now().Add(15 * time.Second))
	for {
		var msg struct {
			Type  string      `json:"type"`
			State model.State `json:"state"`
			Event model.Event `json:"event"`
		}
		if err := c.ReadJSON(&msg); err != nil {
			t.Fatal("waiting revision", want, err)
		}
		if msg.Event.Revision >= want || msg.Type == "snapshot" && msg.State.Revision >= want {
			return
		}
	}
}
func wait(t *testing.T, timeout time.Duration, check func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if check() {
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
	t.Fatal("condition timed out after", timeout)
}
func compose(t *testing.T, args ...string) {
	t.Helper()
	project := os.Getenv("TEST_COMPOSE_PROJECT")
	if project != "syncframe" && project != "syncframe-verify" {
		t.Fatal("fault tests only allow explicitly owned syncframe projects")
	}
	argv := append([]string{"compose", "-p", project, "--profile", "multi"}, args...)
	cmd := exec.Command("docker", argv...)
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("docker %v: %v\n%s", args, err, out)
	}
	t.Logf("docker compose %s", strings.Join(args, " "))
}
func TestAPISecurityAndWebSocketRecovery(t *testing.T) {
	base, c := config(t)
	st, token := create(t, c, base)
	path := base + "/api/boards/" + st.ID
	request(t, &http.Client{Timeout: 5 * time.Second}, "GET", path, nil, 401, nil)
	request(t, &http.Client{Timeout: 5 * time.Second}, "POST", path+"/share", map[string]any{}, 401, nil)
	var shared struct {
		Token string `json:"token"`
	}
	request(t, c, "POST", path+"/share", map[string]any{}, 200, &shared)
	if shared.Token != token {
		t.Fatal("authorized sharing did not preserve board capability")
	}
	req, _ := http.NewRequest("POST", path+"/operations", strings.NewReader(`{}`))
	req.Header.Set("Origin", "https://attacker.invalid")
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatal("cross origin write accepted")
	}
	req, _ = http.NewRequest("GET", path, nil)
	h := http.Header{"Origin": []string{"https://attacker.invalid"}}
	for _, cookie := range c.Jar.Cookies(req.URL) {
		h.Add("Cookie", cookie.String())
	}
	_, resp, err = websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(path, "http")+"/ws?name=Bad&after=0", h)
	if err == nil || resp == nil || resp.StatusCode != 403 {
		t.Fatal("cross-origin WebSocket accepted")
	}
	resp.Body.Close()
	conn := ws(t, c, base, st.ID, 0)
	op := operation(1)
	if err = conn.WriteJSON(map[string]any{"type": "operation", "operation": op}); err != nil {
		t.Fatal(err)
	}
	waitRevision(t, conn, 1)
	conn.Close()
	edit(t, c, base, st.ID, operation(2))
	edit(t, c, base, st.ID, operation(3))
	conn = ws(t, c, base, st.ID, 1)
	defer conn.Close()
	waitRevision(t, conn, 3)
	var state model.State
	request(t, c, "GET", path, nil, 200, &state)
	if state.Revision != 3 || len(state.Objects) != 3 {
		t.Fatal("disconnect catchup failed")
	}
	oversized := operation(9)
	oversized.Object.Text = strings.Repeat("x", 20000)
	request(t, c, "POST", path+"/operations", oversized, 400, nil)
	huge := ws(t, c, base, st.ID, 3)
	defer huge.Close()
	huge.SetReadDeadline(time.Now().Add(5 * time.Second))
	_ = huge.WriteJSON(map[string]any{"type": "operation", "operation": oversized})
	for {
		_, _, err = huge.ReadMessage()
		if err != nil {
			break
		}
	}
	if !websocket.IsCloseError(err, websocket.CloseMessageTooBig) {
		t.Fatalf("expected size rejection, got %v", err)
	}
}
func TestMultiServerAndInfrastructureFaults(t *testing.T) {
	base, c := config(t)
	if os.Getenv("TEST_DOCKER_FAULTS") != "1" {
		t.Skip("TEST_DOCKER_FAULTS=1 required")
	}
	second := os.Getenv("TEST_SECOND_URL")
	if second == "" {
		t.Fatal("TEST_SECOND_URL required")
	}
	s, err := store.Open(os.Getenv("TEST_MYSQL_DSN"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.DB.Close()
	ctx := context.Background()
	st, token := create(t, c, base)
	edit(t, c, base, st.ID, operation(1))
	request(t, c, "POST", second+"/api/boards/"+st.ID+"/session", map[string]string{"token": token}, 200, nil)
	observer := ws(t, c, second, st.ID, 1)
	defer observer.Close()
	edit(t, c, base, st.ID, operation(2))
	waitRevision(t, observer, 2)
	t.Log("SCENARIO 11 PASS: independent servers propagated revision 2")
	compose(t, "stop", "redis")
	defer compose(t, "start", "redis")
	edit(t, c, base, st.ID, operation(3))
	waitRevision(t, observer, 3)
	// Explicitly verify presence degradation, not just SQL persistence.
	observer.SetReadDeadline(time.Now().Add(12 * time.Second))
	for {
		var msg struct {
			Type           string `json:"type"`
			RedisAvailable bool   `json:"redisAvailable"`
		}
		if err = observer.ReadJSON(&msg); err != nil {
			t.Fatal(err)
		}
		if msg.Type == "presence" && !msg.RedisAvailable {
			break
		}
	}
	compose(t, "start", "redis")
	observer.SetReadDeadline(time.Now().Add(18 * time.Second))
	for {
		var msg struct {
			Type           string `json:"type"`
			RedisAvailable bool   `json:"redisAvailable"`
		}
		if err = observer.ReadJSON(&msg); err != nil {
			t.Fatal(err)
		}
		if msg.Type == "presence" && msg.RedisAvailable {
			break
		}
	}
	t.Log("SCENARIO 7 PASS: Redis outage, cross-server polling fallback, presence recovery")
	var before model.State
	request(t, c, "GET", base+"/api/boards/"+st.ID, nil, 200, &before)
	compose(t, "restart", "api")
	wait(t, 40*time.Second, func() bool {
		resp, err := c.Get(base + "/readyz")
		if err != nil {
			return false
		}
		resp.Body.Close()
		return resp.StatusCode == 200
	})
	var restored model.State
	request(t, c, "GET", base+"/api/boards/"+st.ID, nil, 200, &restored)
	if !reflect.DeepEqual(before, restored) {
		t.Fatal("server restart state mismatch")
	}
	reconnected := ws(t, c, base, st.ID, 1)
	waitRevision(t, reconnected, 3)
	reconnected.Close()
	t.Log("SCENARIO 6 PASS: real Go container restarted; state, capability and missing history restored")
	compose(t, "stop", "kafka")
	defer compose(t, "start", "kafka")
	start := time.Now()
	last := edit(t, c, base, st.ID, operation(4))
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatal("edit waited for Kafka", elapsed)
	}
	var pending int
	err = s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM outbox WHERE board_id=? AND published_at IS NULL", st.ID).Scan(&pending)
	if err != nil || pending < 1 {
		t.Fatal("outage outbox not pending", err)
	}
	compose(t, "start", "kafka")
	wait(t, 90*time.Second, func() bool {
		var count int
		err := s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_events WHERE board_id=?", st.ID).Scan(&count)
		return err == nil && count == 4
	})
	wait(t, 30*time.Second, func() bool {
		var count int
		err := s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM outbox WHERE board_id=? AND published_at IS NULL", st.ID).Scan(&count)
		return err == nil && count == 0
	})
	// Real Kafka duplicate delivery, as happens after a publisher's ambiguous success.
	writer := pipeline.Writer(os.Getenv("TEST_KAFKA_BROKER"))
	payload, _ := json.Marshal(last)
	call, cancel := context.WithTimeout(ctx, 10*time.Second)
	err = writer.WriteMessages(call, kafka.Message{Key: []byte(st.ID), Value: payload}, kafka.Message{Key: []byte(st.ID), Value: payload})
	cancel()
	writer.Close()
	if err != nil {
		t.Fatal(err)
	}
	// Reader observes the same tail on a separate group; app audit SQL remains unique.
	reader := kafka.NewReader(kafka.ReaderConfig{Brokers: []string{os.Getenv("TEST_KAFKA_BROKER")}, Topic: pipeline.Topic, Partition: 0, MinBytes: 1, MaxBytes: 1e6})
	defer reader.Close()
	if err = reader.SetOffset(kafka.FirstOffset); err != nil {
		t.Fatal(err)
	}
	copies := 0
	call, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	for copies < 3 {
		m, err := reader.ReadMessage(call)
		if err != nil {
			t.Fatal(err)
		}
		var e model.Event
		_ = json.Unmarshal(m.Value, &e)
		if e.BoardID == st.ID && e.Revision == last.Revision {
			copies++
		}
	}
	// Wait for the application's committed group offset to pass all injected duplicates.
	wait(t, 30*time.Second, func() bool {
		conn, err := kafka.DialContext(ctx, "tcp", os.Getenv("TEST_KAFKA_BROKER"))
		if err != nil {
			return false
		}
		defer conn.Close()
		tail, err := conn.ReadPartitions(pipeline.Topic)
		if err != nil || len(tail) == 0 {
			return false
		}
		broker, err := kafka.DialLeader(ctx, "tcp", os.Getenv("TEST_KAFKA_BROKER"), pipeline.Topic, 0)
		if err != nil {
			return false
		}
		end, err := broker.ReadLastOffset()
		broker.Close()
		if err != nil {
			return false
		}
		client := kafka.Client{Addr: kafka.TCP(os.Getenv("TEST_KAFKA_BROKER"))}
		offsets, err := client.OffsetFetch(ctx, &kafka.OffsetFetchRequest{GroupID: "syncframe-audit-v1", Topics: map[string][]int{pipeline.Topic: {0}}})
		if err != nil {
			return false
		}
		parts := offsets.Topics[pipeline.Topic]
		return len(parts) > 0 && parts[0].CommittedOffset >= end
	})
	var count int
	_ = s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_events WHERE board_id=?", st.ID).Scan(&count)
	if count != 4 {
		t.Fatal("Kafka duplicates changed audit count", count)
	}
	t.Log("SCENARIO 8 PASS: broker outage queued outbox, restart delivered all four events, two real duplicates consumed once")
}
