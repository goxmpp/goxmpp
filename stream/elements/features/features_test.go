package features_test

import (
	"testing"

	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type TestState struct {
	name string
}

func (self *TestState) Name() string {
	return self.name
}

func (self *TestState) SetName(value string) {
	self.name = value
}

func TestSimple(t *testing.T) {
	state := features.NewState()
	ts_w := TestState{}
	ts_w.SetName("test")

	if ts_w.Name() != "test" {
		t.Fatal("Test basic setter/getter failure.")
	}

	state.Push(&ts_w)

	var ts_r *TestState

	if err := state.Get(ts_r); err == nil {
		t.Fatal("Should fail.")
	}

	if err := state.Get(&ts_r); err != nil {
		t.Fatal(err)
	}

	if ts_r.Name() != "test" {
		t.Fatal("Test state getter failure.")
	}
	ts_r.SetName("test2")
	if ts_r.Name() != "test2" {
		t.Fatal("Test state setter failure.")
	}
}
