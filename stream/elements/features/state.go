package features

import (
	"container/list"
	"errors"
	"reflect"
)

type reflectedState struct {
	State interface{}
	Type  reflect.Type
}

type State struct {
	states *list.List
}

func NewState() *State {
	return &State{states: &list.List{}}
}

func getType(val interface{}) reflect.Type {
	vo := reflect.TypeOf(val)
	if vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
	}
	return vo
}

func (self *State) Push(states ...interface{}) {
	for _, state := range states {
		self.states.PushBack(reflectedState{State: state, Type: getType(state)})
	}
}

func (self *State) Get(res interface{}) error {
	// Type of pointer to pointer to value
	typ := reflect.TypeOf(res).Elem()
	if typ.Kind() != reflect.Ptr {
		return errors.New("Should receive a pointer to pointer.")
	}
	typ = typ.Elem()

	// Pointer to pointer to value
	v := reflect.ValueOf(res).Elem()

	for e := self.states.Front(); e != nil; e = e.Next() {
		if state, ok := e.Value.(reflectedState); ok && state.Type == typ {
			v.Set(reflect.ValueOf(state.State))
			return nil
		}
	}

	return errors.New("No state found")
}
