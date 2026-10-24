package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"syncframe/internal/model"
	"syncframe/internal/realtime"
	"syncframe/internal/store"
)

type Server struct {
	Store   *store.Store
	Bus     *realtime.Bus
	Origin  string
	WebDir  string
	mu      sync.Mutex
	sockets map[*websocket.Conn]bool
	closing bool
	wsWait  sync.WaitGroup
}

func New(s *store.Store, b *realtime.Bus, origin, webDir string) *Server {
	return &Server{Store: s, Bus: b, Origin: origin, WebDir: webDir, sockets: map[*websocket.Conn]bool{}}
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { respond(w, 200, map[string]string{"status": "alive"}) })
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()
		if s.Store.DB.PingContext(ctx) != nil {
			respond(w, 503, map[string]string{"status": "mysql unavailable"})
			return
		}
		respond(w, 200, map[string]string{"status": "ready"})
	})
	mux.HandleFunc("POST /api/boards", s.create)
	mux.HandleFunc("POST /api/boards/{id}/session", s.session)
	mux.HandleFunc("POST /api/boards/{id}/share", s.authorize(s.share))
	mux.HandleFunc("GET /api/boards/{id}", s.authorize(s.state))
	mux.HandleFunc("GET /api/boards/{id}/events", s.authorize(s.events))
	mux.HandleFunc("POST /api/boards/{id}/operations", s.authorize(s.operation))
	mux.HandleFunc("GET /api/boards/{id}/ws", s.authorize(s.socket))
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		respond(w, 404, map[string]string{"error": "unknown API route"})
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "HEAD" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		path := filepath.Join(s.WebDir, filepath.Clean("/"+r.URL.Path))
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			http.ServeFile(w, r, path)
			return
		}
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, filepath.Join(s.WebDir, "index.html"))
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; connect-src 'self' ws: wss:; img-src 'self' data:; object-src 'none'; base-uri 'none'; frame-ancestors 'none'")
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Cache-Control", "no-store")
		}
		if r.Method != "GET" && r.Method != "HEAD" {
			if o := r.Header.Get("Origin"); o != "" && !s.allowedOrigin(r, o) {
				respond(w, 403, map[string]string{"error": "origin denied"})
				return
			}
		}
		mux.ServeHTTP(w, r)
	})
}
func (s *Server) allowedOrigin(r *http.Request, origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if s.Origin != "" {
		return origin == s.Origin
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return u.Scheme == scheme && u.Host == r.Host
}
func respond(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func decode(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 16384)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		return model.Invalid(err.Error())
	}
	if err := d.Decode(&struct{}{}); err != io.EOF {
		return model.Invalid("one JSON value required")
	}
	return nil
}
func failure(w http.ResponseWriter, err error) {
	code := 503
	if errors.Is(err, model.ErrInvalid) {
		code = 400
	}
	if errors.Is(err, model.ErrNotFound) {
		code = 404
	}
	if errors.Is(err, model.ErrUnauthorized) {
		code = 401
	}
	if errors.Is(err, model.ErrConflict) {
		code = 409
	}
	if code == 503 {
		slog.Error("request storage failure", "error", err)
		respond(w, 503, map[string]string{"error": "storage temporarily unavailable"})
		return
	}
	respond(w, code, map[string]string{"error": err.Error()})
}
func cookieName(id string) string { return "sf_" + id }
func (s *Server) setCookie(w http.ResponseWriter, r *http.Request, id, token string) {
	http.SetCookie(w, &http.Cookie{Name: cookieName(id), Value: token, Path: "/api/boards/" + id, HttpOnly: true, Secure: r.TLS != nil, SameSite: http.SameSiteStrictMode, MaxAge: 86400 * 30})
}
func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title string `json:"title"`
	}
	if err := decode(w, r, &in); err != nil {
		failure(w, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	st, token, err := s.Store.Create(ctx, in.Title)
	if err != nil {
		failure(w, err)
		return
	}
	s.setCookie(w, r, st.ID, token)
	respond(w, 201, map[string]any{"board": st, "token": token})
}
func (s *Server) session(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Token string `json:"token"`
	}
	if err := decode(w, r, &in); err != nil {
		failure(w, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := s.Store.Auth(ctx, r.PathValue("id"), in.Token); err != nil {
		failure(w, err)
		return
	}
	s.setCookie(w, r, r.PathValue("id"), in.Token)
	respond(w, 200, map[string]bool{"ok": true})
}
func (s *Server) authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName(r.PathValue("id")))
		if err != nil {
			failure(w, model.ErrUnauthorized)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		err = s.Store.Auth(ctx, r.PathValue("id"), c.Value)
		cancel()
		if err != nil {
			failure(w, err)
			return
		}
		next(w, r)
	}
}

// Reopening a board in a new tab retains its HttpOnly cookie, but not its
// sessionStorage. Reveal the already-authorized capability only on explicit,
// same-origin share requests; no new grant or key is created.
func (s *Server) share(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName(r.PathValue("id")))
	if err != nil {
		failure(w, model.ErrUnauthorized)
		return
	}
	respond(w, http.StatusOK, map[string]string{"token": cookie.Value})
}
func (s *Server) state(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var st model.State
	var err error
	if raw := r.URL.Query().Get("revision"); raw != "" {
		n, e := strconv.ParseInt(raw, 10, 64)
		if e != nil {
			failure(w, model.Invalid("invalid revision"))
			return
		}
		st, err = s.Store.At(ctx, r.PathValue("id"), n)
	} else {
		st, err = s.Store.State(ctx, r.PathValue("id"))
	}
	if err != nil {
		failure(w, err)
		return
	}
	respond(w, 200, st)
}
func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.ParseInt(defaultString(r.URL.Query().Get("after"), "0"), 10, 64)
	if err != nil || n < 0 {
		failure(w, model.Invalid("invalid after revision"))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	es, err := s.Store.Events(ctx, r.PathValue("id"), n, 500)
	if err != nil {
		failure(w, err)
		return
	}
	respond(w, 200, es)
}
func (s *Server) operation(w http.ResponseWriter, r *http.Request) {
	var op model.Operation
	if err := decode(w, r, &op); err != nil {
		failure(w, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	e, duplicate, err := s.Store.Apply(ctx, r.PathValue("id"), op)
	if err != nil {
		failure(w, err)
		return
	}
	go s.hint(e.BoardID)
	respond(w, 200, map[string]any{"event": e, "duplicate": duplicate})
}
func (s *Server) hint(board string) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	s.Bus.Notify(ctx, realtime.Signal{BoardID: board, Type: "edit"})
}
func defaultString(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

type inbound struct {
	Type      string           `json:"type"`
	Operation *model.Operation `json:"operation,omitempty"`
	X         float64          `json:"x,omitempty"`
	Y         float64          `json:"y,omitempty"`
}

func (s *Server) socket(w http.ResponseWriter, r *http.Request) {
	after, err := strconv.ParseInt(defaultString(r.URL.Query().Get("after"), "0"), 10, 64)
	if err != nil || after < 0 {
		failure(w, model.Invalid("invalid revision"))
		return
	}
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	if name == "" || len(name) > 60 {
		failure(w, model.Invalid("name must be 1-60 bytes"))
		return
	}
	call, stop := context.WithTimeout(r.Context(), 3*time.Second)
	current, stateErr := s.Store.State(call, r.PathValue("id"))
	stop()
	if stateErr != nil {
		failure(w, stateErr)
		return
	}
	if after > current.Revision {
		failure(w, model.Invalid("resume revision is ahead of the server"))
		return
	}
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { o := r.Header.Get("Origin"); return o == "" || s.allowedOrigin(r, o) }, HandshakeTimeout: 3 * time.Second}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	s.mu.Lock()
	if s.closing || len(s.sockets) >= 200 {
		s.mu.Unlock()
		conn.Close()
		return
	}
	s.sockets[conn] = true
	s.wsWait.Add(1)
	s.mu.Unlock()
	defer s.wsWait.Done()
	defer func() { s.mu.Lock(); delete(s.sockets, conn); s.mu.Unlock(); conn.Close() }()
	p := s.Bus.Add(r.PathValue("id"), name)
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	defer func() {
		c, stop := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer stop()
		s.Bus.Remove(c, p)
	}()
	outgoing := make(chan any, 32)
	readDone := make(chan struct{})
	conn.SetReadLimit(16384)
	_ = conn.SetReadDeadline(time.Now().Add(40 * time.Second))
	conn.SetPongHandler(func(string) error { return conn.SetReadDeadline(time.Now().Add(40 * time.Second)) })
	go func() {
		defer close(readDone)
		defer cancel()
		lastCursor := time.Time{}
		window := time.Now()
		edits := 0
		for {
			var in inbound
			if err := conn.ReadJSON(&in); err != nil {
				return
			}
			switch in.Type {
			case "operation":
				if time.Since(window) > time.Second {
					window = time.Now()
					edits = 0
				}
				edits++
				if edits > 60 {
					return
				}
				if in.Operation == nil {
					return
				}
				in.Operation.Actor = name
				call, stop := context.WithTimeout(ctx, 5*time.Second)
				e, dup, err := s.Store.Apply(call, p.Board, *in.Operation)
				stop()
				var message any
				if err != nil {
					message = map[string]any{"type": "rejected", "operationId": in.Operation.ID, "error": "operation failed: " + publicError(err), "retryable": !errors.Is(err, model.ErrConflict) && publicError(err) == "storage temporarily unavailable; retry the same operation ID"}
				} else {
					message = map[string]any{"type": "ack", "operationId": e.Operation.ID, "revision": e.Revision, "duplicate": dup}
					go s.hint(p.Board)
				}
				select {
				case outgoing <- message:
				case <-ctx.Done():
					return
				default:
					return
				}
			case "cursor":
				if time.Since(lastCursor) < 40*time.Millisecond {
					continue
				}
				if math.IsNaN(in.X) || math.IsNaN(in.Y) || math.Abs(in.X) > 10000 || math.Abs(in.Y) > 10000 {
					continue
				}
				lastCursor = time.Now()
				call, stop := context.WithTimeout(ctx, 150*time.Millisecond)
				s.Bus.Notify(call, realtime.Signal{BoardID: p.Board, Type: "cursor", Peer: &p.User, X: in.X, Y: in.Y})
				stop()
			default:
				return
			}
		}
	}()
	defer func() { cancel(); conn.Close(); <-readDone }()
	write := func(v any) error {
		_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		return conn.WriteJSON(v)
	}
	last := after
	syncEvents := func() error {
		call, stop := context.WithTimeout(ctx, 5*time.Second)
		defer stop()
		if last == 0 {
			st, err := s.Store.State(call, p.Board)
			if err != nil {
				return err
			}
			if err = write(map[string]any{"type": "snapshot", "state": st, "selfId": p.User.ID}); err != nil {
				return err
			}
			last = st.Revision
			return nil
		}
		es, err := s.Store.Events(call, p.Board, last, 500)
		if err != nil {
			return err
		}
		for _, e := range es {
			if err = write(map[string]any{"type": "event", "event": e}); err != nil {
				return err
			}
			last = e.Revision
		}
		return nil
	}
	if err = write(map[string]any{"type": "hello", "selfId": p.User.ID}); err != nil {
		return
	}
	if err = syncEvents(); err != nil {
		return
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	heartbeat := time.NewTicker(5 * time.Second)
	defer heartbeat.Stop()
	presence := func() error {
		call, stop := context.WithTimeout(ctx, 1200*time.Millisecond)
		defer stop()
		users, redisOK := s.Bus.Presence(call, p)
		return write(map[string]any{"type": "presence", "users": users, "redisAvailable": redisOK})
	}
	if presence() != nil {
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.Wake:
			if syncEvents() != nil {
				return
			}
		case <-ticker.C:
			if syncEvents() != nil {
				return
			}
		case <-heartbeat.C:
			if presence() != nil {
				return
			}
			if err = conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second)); err != nil {
				return
			}
		case message := <-outgoing:
			if write(message) != nil {
				return
			}
		case signal := <-p.Ephemeral:
			if write(signal) != nil {
				return
			}
		}
	}
}
func publicError(err error) string {
	if errors.Is(err, model.ErrConflict) || errors.Is(err, model.ErrInvalid) {
		return err.Error()
	}
	return "storage temporarily unavailable; retry the same operation ID"
}
func (s *Server) Close() {
	s.mu.Lock()
	s.closing = true
	for c := range s.sockets {
		c.Close()
	}
	s.mu.Unlock()
	s.wsWait.Wait()
}
