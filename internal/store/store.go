package store

import (
	"context"
	"crypto/subtle"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"syncframe/internal/model"
)

//go:embed schema.sql
var schema string

type Store struct{ DB *sql.DB }

func Open(dsn string) (*Store, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(16)
	db.SetMaxIdleConns(8)
	db.SetConnMaxLifetime(3 * time.Minute)
	return &Store{DB: db}, nil
}
func (s *Store) Migrate(ctx context.Context) error {
	for _, q := range strings.Split(schema, ";") {
		if strings.TrimSpace(q) == "" {
			continue
		}
		if _, err := s.DB.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}
func (s *Store) Create(ctx context.Context, title string) (model.State, string, error) {
	title = strings.TrimSpace(title)
	if title == "" || len(title) > 120 {
		return model.State{}, "", model.Invalid("title must be 1-120 bytes")
	}
	st := model.State{ID: model.ID(), Title: title, Objects: map[string]model.Object{}}
	token := model.ID() + model.ID()
	b, _ := json.Marshal(st)
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return st, "", err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, "INSERT INTO boards(id,title,token_hash,state) VALUES(?,?,?,?)", st.ID, title, model.Hash(token), b); err != nil {
		return st, "", err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO snapshots(board_id,revision,state) VALUES(?,0,?)", st.ID, b); err != nil {
		return st, "", err
	}
	err = tx.Commit()
	return st, token, err
}
func (s *Store) Auth(ctx context.Context, id, token string) error {
	if !model.ValidID(id) || len(token) != 64 {
		return model.ErrUnauthorized
	}
	var h string
	err := s.DB.QueryRowContext(ctx, "SELECT token_hash FROM boards WHERE id=?", id).Scan(&h)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrUnauthorized
	}
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare([]byte(h), []byte(model.Hash(token))) != 1 {
		return model.ErrUnauthorized
	}
	return nil
}
func (s *Store) State(ctx context.Context, id string) (model.State, error) {
	var b []byte
	err := s.DB.QueryRowContext(ctx, "SELECT state FROM boards WHERE id=?", id).Scan(&b)
	if errors.Is(err, sql.ErrNoRows) {
		err = model.ErrNotFound
	}
	var st model.State
	if err == nil {
		err = json.Unmarshal(b, &st)
	}
	return st, err
}
func (s *Store) Apply(ctx context.Context, id string, op model.Operation) (model.Event, bool, error) {
	if err := model.Validate(op); err != nil {
		return model.Event{}, false, err
	}
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return model.Event{}, false, err
	}
	defer tx.Rollback()
	var b []byte
	if err = tx.QueryRowContext(ctx, "SELECT state FROM boards WHERE id=? FOR UPDATE", id).Scan(&b); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = model.ErrNotFound
		}
		return model.Event{}, false, err
	}
	var st model.State
	if err = json.Unmarshal(b, &st); err != nil {
		return model.Event{}, false, err
	}
	var old []byte
	var fingerprint string
	err = tx.QueryRowContext(ctx, "SELECT payload,fingerprint FROM events WHERE board_id=? AND operation_id=?", id, op.ID).Scan(&old, &fingerprint)
	if err == nil {
		if fingerprint != model.Fingerprint(op) {
			return model.Event{}, false, fmt.Errorf("%w: operation ID reused with different payload", model.ErrConflict)
		}
		var e model.Event
		err = json.Unmarshal(old, &e)
		return e, true, err
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return model.Event{}, false, err
	}
	if st.Revision >= 100000 {
		return model.Event{}, false, model.Invalid("board event limit reached")
	}
	obj, err := model.Apply(st, op)
	if err != nil {
		return model.Event{}, false, err
	}
	e := model.Event{BoardID: id, Revision: st.Revision + 1, Operation: op, Object: obj, At: time.Now().UTC()}
	st.Revision = e.Revision
	st.Objects[obj.ID] = obj
	eb, _ := json.Marshal(e)
	sb, _ := json.Marshal(st)
	if _, err = tx.ExecContext(ctx, "INSERT INTO events(board_id,revision,operation_id,fingerprint,payload) VALUES(?,?,?,?,?)", id, e.Revision, op.ID, model.Fingerprint(op), eb); err != nil {
		return e, false, err
	}
	if _, err = tx.ExecContext(ctx, "UPDATE boards SET revision=?,state=? WHERE id=?", e.Revision, sb, id); err != nil {
		return e, false, err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO outbox(board_id,revision,payload) VALUES(?,?,?)", id, e.Revision, eb); err != nil {
		return e, false, err
	}
	if e.Revision%50 == 0 {
		if _, err = tx.ExecContext(ctx, "INSERT INTO snapshots(board_id,revision,state) VALUES(?,?,?)", id, e.Revision, sb); err != nil {
			return e, false, err
		}
	}
	return e, false, tx.Commit()
}
func (s *Store) Events(ctx context.Context, id string, after int64, limit int) ([]model.Event, error) {
	if limit < 1 || limit > 1000 {
		limit = 500
	}
	rows, err := s.DB.QueryContext(ctx, "SELECT payload FROM events WHERE board_id=? AND revision>? ORDER BY revision LIMIT ?", id, after, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Event{}
	for rows.Next() {
		var b []byte
		if err = rows.Scan(&b); err != nil {
			return nil, err
		}
		var e model.Event
		if err = json.Unmarshal(b, &e); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
func (s *Store) At(ctx context.Context, id string, revision int64) (model.State, error) {
	current, err := s.State(ctx, id)
	if err != nil {
		return current, err
	}
	if revision < 0 || revision > current.Revision {
		return model.State{}, model.Invalid("revision out of range")
	}
	var b []byte
	err = s.DB.QueryRowContext(ctx, "SELECT state FROM snapshots WHERE board_id=? AND revision<=? ORDER BY revision DESC LIMIT 1", id, revision).Scan(&b)
	if err != nil {
		return model.State{}, err
	}
	var st model.State
	if err = json.Unmarshal(b, &st); err != nil {
		return st, err
	}
	for st.Revision < revision {
		events, err := s.Events(ctx, id, st.Revision, 1000)
		if err != nil {
			return st, err
		}
		if len(events) == 0 {
			return st, errors.New("incomplete history")
		}
		for _, e := range events {
			if e.Revision > revision {
				break
			}
			if err = model.Replay(&st, e); err != nil {
				return st, err
			}
		}
	}
	return st, nil
}
func (s *Store) Audit(ctx context.Context, e model.Event) error {
	b, _ := json.Marshal(e)
	_, err := s.DB.ExecContext(ctx, "INSERT INTO audit_events(board_id,revision,payload) VALUES(?,?,?) ON DUPLICATE KEY UPDATE revision=revision", e.BoardID, e.Revision, b)
	return err
}
