package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"syncframe/internal/model"
	"testing"
	"time"
)

func integration(t *testing.T) *Store {
	t.Helper()
	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("TEST_MYSQL_DSN not set: real MySQL integration skipped")
	}
	s, err := Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.DB.Close() })
	ctx, c := context.WithTimeout(context.Background(), 10*time.Second)
	defer c()
	if err = s.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	return s
}
func makeOp(i int) model.Operation {
	id := fmt.Sprintf("object-%d", i)
	o := model.Object{ID: id, Kind: "rect", X: float64(i), Y: 20, Width: 120, Height: 80, Color: "#dce7ff"}
	return model.Operation{ID: model.ID(), Actor: fmt.Sprintf("user-%d", i), Kind: "create", ObjectID: id, Object: &o}
}
func TestMySQLConcurrencyAndReplay(t *testing.T) {
	s := integration(t)
	ctx := context.Background()
	st, _, err := s.Create(ctx, "20 concurrent creators")
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	errs := make(chan error, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) { defer wg.Done(); _, _, err := s.Apply(ctx, st.ID, makeOp(i)); errs <- err }(i)
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	latest, err := s.State(ctx, st.ID)
	if err != nil || latest.Revision != 20 || len(latest.Objects) != 20 {
		t.Fatalf("creations: %+v %v", latest, err)
	}
	events, err := s.Events(ctx, st.ID, 0, 1000)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for i, e := range events {
		if e.Revision != int64(i+1) || seen[e.Operation.ID] {
			t.Fatal("revision or operation uniqueness")
		}
		seen[e.Operation.ID] = true
	}
	// Same object, competing field updates: final x must equal the highest committed revision.
	var maxRevision int64
	var expected float64
	var mu sync.Mutex
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			x := float64(i * 10)
			op := model.Operation{ID: model.ID(), Actor: "writer", Kind: "update", ObjectID: "object-0", BaseRevision: 20, Patch: &model.Patch{X: &x}}
			e, _, err := s.Apply(ctx, st.ID, op)
			if err != nil {
				t.Error(err)
				return
			}
			mu.Lock()
			if e.Revision > maxRevision {
				maxRevision = e.Revision
				expected = x
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	latest, err = s.State(ctx, st.ID)
	if err != nil || latest.Objects["object-0"].X != expected || latest.Revision != 40 {
		t.Fatal("last committed field write mismatch", err)
	}
	// Exercise snapshot at revision 50 and suffix reconstruction, not only revision zero.
	for i := 20; i < 35; i++ {
		if _, _, err = s.Apply(ctx, st.ID, makeOp(i)); err != nil {
			t.Fatal(err)
		}
	}
	all, err := s.Events(ctx, st.ID, 0, 1000)
	if err != nil {
		t.Fatal(err)
	}
	replay := model.State{ID: st.ID, Title: st.Title, Objects: map[string]model.Object{}}
	for _, e := range all {
		if err = model.Replay(&replay, e); err != nil {
			t.Fatal(err)
		}
	}
	latest, _ = s.State(ctx, st.ID)
	if !reflect.DeepEqual(replay, latest) {
		t.Fatal("replay != latest")
	}
	at, err := s.At(ctx, st.ID, 53)
	if err != nil || at.Revision != 53 || len(at.Objects) != 33 {
		t.Fatalf("snapshot suffix: %+v %v", at, err)
	}
	if after, _ := s.State(ctx, st.ID); after.Revision != 55 {
		t.Fatal("history modified live document")
	}
	// Reopen independent DB pool to prove state is not held in Go memory.
	second, err := Open(os.Getenv("TEST_MYSQL_DSN"))
	if err != nil {
		t.Fatal(err)
	}
	defer second.DB.Close()
	restored, err := second.State(ctx, st.ID)
	if err != nil || !reflect.DeepEqual(restored, latest) {
		t.Fatal("reopened store mismatch", err)
	}
}
func TestMySQLDuplicateAndDeletion(t *testing.T) {
	s := integration(t)
	ctx := context.Background()
	st, token, err := s.Create(ctx, "duplicates")
	if err != nil {
		t.Fatal(err)
	}
	if s.Auth(ctx, st.ID, token) != nil || !errors.Is(s.Auth(ctx, st.ID, "bad"), model.ErrUnauthorized) {
		t.Fatal("capability authentication")
	}
	op := makeOp(1)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e, _, err := s.Apply(ctx, st.ID, op)
			if err != nil || e.Revision != 1 {
				t.Error("duplicate", err)
			}
		}()
	}
	wg.Wait()
	latest, _ := s.State(ctx, st.ID)
	if latest.Revision != 1 {
		t.Fatal("duplicate revision")
	}
	changed := op
	changed.Actor = "different"
	if _, _, err = s.Apply(ctx, st.ID, changed); !errors.Is(err, model.ErrConflict) {
		t.Fatal("operation ID payload conflict not detected")
	}
	deletion := model.Operation{ID: model.ID(), Actor: "delete", Kind: "delete", ObjectID: op.ObjectID}
	if _, _, err = s.Apply(ctx, st.ID, deletion); err != nil {
		t.Fatal(err)
	}
	x := 70.0
	update := model.Operation{ID: model.ID(), Actor: "stale", Kind: "update", ObjectID: op.ObjectID, BaseRevision: 1, Patch: &model.Patch{X: &x}}
	if _, _, err = s.Apply(ctx, st.ID, update); !errors.Is(err, model.ErrConflict) {
		t.Fatal("resurrection accepted")
	}
	if _, _, err = s.Apply(ctx, st.ID, op); err != nil {
		t.Fatal("original duplicate after deletion must ack", err)
	}
	latest, _ = s.State(ctx, st.ID)
	if !latest.Objects[op.ObjectID].Deleted || latest.Revision != 2 {
		t.Fatal("tombstone lost")
	}
}
func TestMySQLRollbackInjectedFailure(t *testing.T) {
	s := integration(t)
	if os.Getenv("TEST_ALLOW_FAULTS") != "1" {
		t.Skip("TEST_ALLOW_FAULTS=1 required for test-only SQL trigger")
	}
	ctx := context.Background()
	st, _, err := s.Create(ctx, "rollback")
	if err != nil {
		t.Fatal(err)
	}
	// Trigger fails the last write after event insertion and board state update.
	trigger := "sf_fault_" + model.ID()
	adminDSN := os.Getenv("TEST_ADMIN_DSN")
	if adminDSN == "" {
		t.Fatal("TEST_ADMIN_DSN required to create a test-only trigger without elevating the app account")
	}
	admin, err := Open(adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.DB.Close()
	q := fmt.Sprintf("CREATE TRIGGER `%s` BEFORE INSERT ON outbox FOR EACH ROW BEGIN IF NEW.board_id='%s' THEN SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT='SYNCFRAME injected failure'; END IF; END", trigger, st.ID)
	if _, err = admin.DB.ExecContext(ctx, q); err != nil {
		t.Fatal(err)
	}
	defer admin.DB.ExecContext(ctx, "DROP TRIGGER IF EXISTS `"+trigger+"`")
	if _, _, err = s.Apply(ctx, st.ID, makeOp(1)); err == nil {
		t.Fatal("fault did not fire")
	}
	latest, _ := s.State(ctx, st.ID)
	if latest.Revision != 0 || len(latest.Objects) != 0 {
		t.Fatal("partial state persisted")
	}
	for _, table := range []string{"events", "outbox"} {
		var count int
		if err = s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table+" WHERE board_id=?", st.ID).Scan(&count); err != nil || count != 0 {
			t.Fatal("half event in "+table, err)
		}
	}
	if _, err = admin.DB.ExecContext(ctx, "DROP TRIGGER `"+trigger+"`"); err != nil {
		t.Fatal(err)
	}
	e, _, err := s.Apply(ctx, st.ID, makeOp(1))
	if err != nil || e.Revision != 1 {
		t.Fatal("revision skipped after rollback", err)
	}
}
func TestAuditIdempotence(t *testing.T) {
	s := integration(t)
	ctx := context.Background()
	st, _, err := s.Create(ctx, "audit")
	if err != nil {
		t.Fatal(err)
	}
	e, _, err := s.Apply(ctx, st.ID, makeOp(1))
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		if err = s.Audit(ctx, e); err != nil {
			t.Fatal(err)
		}
	}
	var count int
	var raw []byte
	err = s.DB.QueryRowContext(ctx, "SELECT COUNT(*),ANY_VALUE(payload) FROM audit_events WHERE board_id=?", st.ID).Scan(&count, &raw)
	if err != nil || count != 1 {
		t.Fatal("audit duplicates", err)
	}
	var got model.Event
	if json.Unmarshal(raw, &got) != nil || got.Revision != e.Revision {
		t.Fatal("audit payload mismatch")
	}
}
