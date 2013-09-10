package xep

import "errors"
import . "github.com/dotdoom/goxmpp/stream/decoder"

const (
	UnknownKey = errors.New("There is no such handler registered")
)

type HandlerRegistrator map[string]decoder.ElementHandler

func NewHandlerRegistrator() Registrator {
	return make(Registrator)
}

func (self HandlerRegistrator) Register(key string, handler decoder.ElementHandler) {
	r := *self
	r[key] = handler
	self = &r
}

func (self HandlerRegistrator) GetHandler(key string) (decoder.ElementHandler, error) {
	if handler, ok := self[key]; ok {
		return handler, nil
	}
	return nil, UnknownKey
}
