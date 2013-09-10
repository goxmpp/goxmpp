package xep

import "errors"
import . "github.com/dotdoom/goxmpp/stream/stanza"

const (
	UnknownKey = errors.New("There is no such handler registered")
)

type Registrator map[string]ElementHandler

func NewRegistrator() Registrator {
	return make(Registrator)
}

func (self *Registrator) Register(key string, handler ElementHandler) {
	r := *self
	r[key] = handler
	self = &r
}

func (self Registrator) GetHandler(key string) (ElementHandler, error) {
	if handler, ok := self[key]; ok {
		return handler, nil
	}
	return nil, UnknownKey
}
