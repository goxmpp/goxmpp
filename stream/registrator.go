package stream

import "errors"

var UnknownKey = errors.New("There is no such handler registered")

type ElementHandlerRegistrator map[string]ElementHandler

func NewElementHandlerRegistrator() ElementHandlerRegistrator {
	return make(ElementHandlerRegistrator)
}

func (self ElementHandlerRegistrator) Register(key string, handler ElementHandler) {
	self[key] = handler
}

func (self ElementHandlerRegistrator) GetHandler(key string) (ElementHandler, error) {
	if handler, ok := self[key]; ok {
		return handler, nil
	}
	return nil, UnknownKey
}
