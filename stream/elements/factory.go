package elements

import "errors"

type Constructor func() Element
type Factory map[string]Constructor

func NewFactory() Factory {
	return make(Factory)
}

func (self Factory) AddConstructor(key string, constructor Constructor) {
	self[key] = constructor
}

func (self Factory) Create(key string) (Element, error) {
	if constructor, ok := self[key]; ok {
		return constructor(), nil
	}

	// This is default constructor if defined
	if constructor, ok := self["*"]; ok {
		return constructor(), nil
	}
	return nil, errors.New("Wrong key " + key)
}
