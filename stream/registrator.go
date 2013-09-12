package stream

import "errors"

var UnknownKey = errors.New("There is no such handler registered")

type ElementGenerator func() Element
type ElementGeneratorRegistrator map[string]ElementGenerator

func NewElementGeneratorRegistrator() ElementGeneratorRegistrator {
	return make(ElementGeneratorRegistrator)
}

func (self ElementGeneratorRegistrator) Register(key string, generator ElementGenerator) {
	self[key] = generator
}

func (self ElementGeneratorRegistrator) GetHandler(key string) (Element, error) {
	if generator, ok := self[key]; ok {
		return generator(), nil
	}
	return nil, errors.New("Wrong key" + key)
}
