package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"syncframe/internal/model"
	"syncframe/internal/realtime"
	"syncframe/internal/store"
)

func TestConcurrentSocketsAndShutdown(t *testing.T) {
	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("TEST_MYSQL_DSN not set")
	}
	s, err := store.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer s.DB.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err = s.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:63790"
	}
	bus := realtime.New(addr)
	defer bus.Redis.Close()
	busDone := make(chan struct{})
	go func() { defer close(busDone); bus.Run(ctx) }()
	defer func() { cancel(); <-busDone }()
	app := New(s, bus, "", "../../web/dist")
	server := httptest.NewServer(app.Handler())
	defer server.Close()
	defer app.Close()
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Timeout: 5 * time.Second}
	response, err := client.Post(server.URL+"/api/boards", "application/json", bytes.NewBufferString(`{"title":"Concurrent socket integration"}`))
	if err != nil {
		t.Fatal(err)
	}
	var created struct {
		Board model.State `json:"board"`
	}
	err = json.NewDecoder(response.Body).Decode(&created)
	response.Body.Close()
	if err != nil || response.StatusCode != 201 {
		t.Fatal("create board", err)
	}
	req, _ := http.NewRequest("GET", server.URL+"/api/boards/"+created.Board.ID, nil)
	headers := http.Header{"Origin": []string{server.URL}}
	for _, c := range jar.Cookies(req.URL) {
		headers.Add("Cookie", c.String())
	}
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"/api/boards/"+created.Board.ID+"/ws?name=Race&after=0", headers)
		if err != nil {
			t.Fatal(err)
		}
		wg.Add(1)
		go func(i int, c *websocket.Conn) {
			defer wg.Done()
			defer c.Close()
			id := fmt.Sprintf("object-%d", i)
			op := model.Operation{ID: model.ID(), Actor: "Race", Kind: "create", ObjectID: id, Object: &model.Object{ID: id, Kind: "note", Width: 200, Height: 150, Color: "#fff2b3"}}
			if err := c.WriteJSON(map[string]any{"type": "operation", "operation": op}); err != nil {
				t.Error(err)
				return
			}
			c.SetReadDeadline(time.Now().Add(15 * time.Second))
			for {
				var msg struct {
					Type  string      `json:"type"`
					State model.State `json:"state"`
					Event model.Event `json:"event"`
				}
				if err := c.ReadJSON(&msg); err != nil {
					t.Error(err)
					return
				}
				if msg.Event.Revision == 5 || msg.Type == "snapshot" && msg.State.Revision == 5 {
					return
				}
			}
		}(i, conn)
	}
	wg.Wait()
	state, err := s.State(ctx, created.Board.ID)
	if err != nil || state.Revision != 5 || len(state.Objects) != 5 {
		t.Fatal("concurrent WebSocket convergence", err)
	}
	done := make(chan struct{})
	go func() { app.Close(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("WebSocket shutdown leaked handlers")
	}
}
