package stream

import "errors"

var UnknownKey = errors.New("There is no such handler registered")

type ElementGenerator func() Element
type ElementHandlerRegistrator map[string]ElementGenerator

func NewElementHandlerRegistrator() ElementHandlerRegistrator {
	return make(ElementHandlerRegistrator)
}

func (self ElementHandlerRegistrator) Register(key string, generator ElementGenerator) {
	self[key] = generator
}

func (self ElementHandlerRegistrator) GetHandler(key string) (Element, error) {
	if generator, ok := self[key]; ok {
		return generator(), nil
	}
	return nil, errors.New("Wrong key" + key)
}
