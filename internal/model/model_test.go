package model

import (
	"errors"
	"testing"
)

func sample() Object {
	return Object{ID: "object-1", Kind: "note", X: 20, Y: 30, Width: 200, Height: 150, Color: "#fff2b3", Text: "idea"}
}
func TestConflictAndTombstone(t *testing.T) {
	s := State{Objects: map[string]Object{}}
	o := sample()
	op := Operation{ID: "op-1", Actor: "A", Kind: "create", ObjectID: o.ID, Object: &o}
	result, err := Apply(s, op)
	if err != nil {
		t.Fatal(err)
	}
	s.Objects[o.ID] = result
	s.Revision = 1
	x := 99.0
	op = Operation{ID: "op-2", Actor: "B", Kind: "update", ObjectID: o.ID, Patch: &Patch{X: &x}}
	result, err = Apply(s, op)
	if err != nil || result.X != 99 || result.Text != "idea" {
		t.Fatalf("field patch: %+v %v", result, err)
	}
	s.Objects[o.ID] = result
	op.Kind = "delete"
	op.Patch = nil
	result, err = Apply(s, op)
	if err != nil {
		t.Fatal(err)
	}
	s.Objects[o.ID] = result
	op.Kind = "update"
	op.Patch = &Patch{X: &x}
	if _, err = Apply(s, op); !errors.Is(err, ErrConflict) {
		t.Fatal("deleted object resurrected", err)
	}
	op.Kind = "create"
	op.Patch = nil
	op.Object = &o
	if _, err = Apply(s, op); !errors.Is(err, ErrConflict) {
		t.Fatal("tombstone ID reused", err)
	}
}
func TestReplayRejectsGapAndIgnoresDuplicate(t *testing.T) {
	s := State{Objects: map[string]Object{}}
	e := Event{Revision: 2, Object: sample()}
	if Replay(&s, e) == nil || s.Revision != 0 {
		t.Fatal("gap applied")
	}
	e.Revision = 1
	if err := Replay(&s, e); err != nil {
		t.Fatal(err)
	}
	if err := Replay(&s, e); err != nil || s.Revision != 1 {
		t.Fatal("duplicate changed revision")
	}
}
func TestValidation(t *testing.T) {
	o := sample()
	op := Operation{ID: "a", Actor: "A", Kind: "create", ObjectID: o.ID, Object: &o}
	o.Width = -1
	if Validate(op) == nil {
		t.Fatal("negative width accepted")
	}
	o = sample()
	o.Color = "url(javascript:x)"
	if Validate(op) == nil {
		t.Fatal("invalid color accepted")
	}
	o = sample()
	op.BaseRevision = 5
	if _, err := Apply(State{Revision: 1, Objects: map[string]Object{}}, op); !errors.Is(err, ErrConflict) {
		t.Fatal("future revision accepted")
	}
}
