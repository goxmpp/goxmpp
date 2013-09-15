package stream

import "errors"

type ElementConstructor func() Element
type ElementFactory map[string]ElementConstructor

func NewElementFactory() ElementFactory {
	return make(ElementFactory)
}

func (self ElementFactory) AddConstructor(key string, constructor ElementConstructor) {
	self[key] = constructor
}

func (self ElementFactory) Create(key string) (Element, error) {
	if constructor, ok := self[key]; ok {
		return constructor(), nil
	}

	// This is default constructor if defined
	if constructor, ok := self["*"]; ok {
		return constructor(), nil
	}
	return nil, errors.New("Wrong key " + key)
}
