package features

import (
	"container/list"
)

type State struct {
	States *list.List
}

func NewState() *State {
	return &State{States: &list.List{}}
}
