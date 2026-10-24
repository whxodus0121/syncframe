package model

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")
var ErrUnauthorized = errors.New("unauthorized")
var ErrInvalid = errors.New("invalid input")

func Invalid(message string) error { return fmt.Errorf("%w: %s", ErrInvalid, message) }

var idPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,80}$`)
var colorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func ID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}
func ValidID(s string) bool { return idPattern.MatchString(s) }
func Hash(s string) string  { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }

type Object struct {
	ID      string  `json:"id"`
	Kind    string  `json:"kind"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
	Color   string  `json:"color"`
	Text    string  `json:"text"`
	Deleted bool    `json:"deleted"`
}
type Patch struct {
	X      *float64 `json:"x,omitempty"`
	Y      *float64 `json:"y,omitempty"`
	Width  *float64 `json:"width,omitempty"`
	Height *float64 `json:"height,omitempty"`
	Color  *string  `json:"color,omitempty"`
	Text   *string  `json:"text,omitempty"`
}
type Operation struct {
	ID           string  `json:"id"`
	Actor        string  `json:"actor"`
	Kind         string  `json:"kind"`
	ObjectID     string  `json:"objectId"`
	BaseRevision int64   `json:"baseRevision"`
	Object       *Object `json:"object,omitempty"`
	Patch        *Patch  `json:"patch,omitempty"`
}
type Event struct {
	BoardID   string    `json:"boardId"`
	Revision  int64     `json:"revision"`
	Operation Operation `json:"operation"`
	Object    Object    `json:"object"`
	At        time.Time `json:"at"`
}
type State struct {
	ID       string            `json:"id"`
	Title    string            `json:"title"`
	Revision int64             `json:"revision"`
	Objects  map[string]Object `json:"objects"`
}

func Fingerprint(op Operation) string { b, _ := json.Marshal(op); return Hash(string(b)) }
func Validate(op Operation) error {
	if err := validateOperation(op); err != nil {
		return Invalid(err.Error())
	}
	return nil
}
func validateOperation(op Operation) error {
	if !ValidID(op.ID) || !ValidID(op.ObjectID) || strings.TrimSpace(op.Actor) == "" || len(op.Actor) > 60 || op.BaseRevision < 0 {
		return errors.New("invalid operation identity, author or revision")
	}
	switch op.Kind {
	case "create":
		if op.Object == nil || op.Patch != nil || op.Object.ID != op.ObjectID || op.Object.Deleted {
			return errors.New("invalid create payload")
		}
		return validateObject(*op.Object)
	case "update":
		if op.Patch == nil || op.Object != nil {
			return errors.New("update requires patch only")
		}
		p := op.Patch
		if p.X == nil && p.Y == nil && p.Width == nil && p.Height == nil && p.Color == nil && p.Text == nil {
			return errors.New("empty patch")
		}
	case "delete":
		if op.Object != nil || op.Patch != nil {
			return errors.New("delete has no object payload")
		}
	default:
		return errors.New("unknown operation kind")
	}
	return nil
}
func ValidateObject(o Object) error {
	if err := validateObject(o); err != nil {
		return Invalid(err.Error())
	}
	return nil
}
func validateObject(o Object) error {
	if !ValidID(o.ID) || (o.Kind != "rect" && o.Kind != "ellipse" && o.Kind != "text" && o.Kind != "note") {
		return errors.New("invalid object identity or kind")
	}
	for _, v := range []float64{o.X, o.Y, o.Width, o.Height} {
		if math.IsNaN(v) || math.IsInf(v, 0) || math.Abs(v) > 10000 {
			return errors.New("coordinate outside canvas limits")
		}
	}
	if o.Width < 20 || o.Height < 20 || o.Width > 2000 || o.Height > 2000 || len(o.Text) > 2000 || !colorPattern.MatchString(o.Color) {
		return errors.New("invalid dimensions, color or text length")
	}
	return nil
}
func Apply(s State, op Operation) (Object, error) {
	if err := Validate(op); err != nil {
		return Object{}, err
	}
	if op.BaseRevision > s.Revision {
		return Object{}, fmt.Errorf("%w: future base revision", ErrConflict)
	}
	current, exists := s.Objects[op.ObjectID]
	switch op.Kind {
	case "create":
		if exists {
			return Object{}, fmt.Errorf("%w: object ID already used", ErrConflict)
		}
		if len(s.Objects) >= 500 {
			return Object{}, Invalid("board object limit reached")
		}
		return *op.Object, nil
	case "update", "delete":
		if !exists || current.Deleted {
			return Object{}, fmt.Errorf("%w: object missing or deleted", ErrConflict)
		}
	}
	if op.Kind == "delete" {
		current.Deleted = true
		return current, nil
	}
	p := op.Patch
	if p.X != nil {
		current.X = *p.X
	}
	if p.Y != nil {
		current.Y = *p.Y
	}
	if p.Width != nil {
		current.Width = *p.Width
	}
	if p.Height != nil {
		current.Height = *p.Height
	}
	if p.Color != nil {
		current.Color = *p.Color
	}
	if p.Text != nil {
		current.Text = *p.Text
	}
	return current, ValidateObject(current)
}
func Replay(s *State, e Event) error {
	if e.Revision <= s.Revision {
		return nil
	}
	if e.Revision != s.Revision+1 {
		return errors.New("revision gap")
	}
	if s.Objects == nil {
		s.Objects = map[string]Object{}
	}
	s.Objects[e.Object.ID] = e.Object
	s.Revision = e.Revision
	return nil
}
