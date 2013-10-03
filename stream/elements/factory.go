package elements

import "errors"

// Create an (empty) Parsable to parse XML into
type Constructor func() Parsable

// Maintain a mapping between tag names (and namespaces) and Constructors
type Factory map[string]Constructor

func (self Factory) AddConstructor(key string, constructor Constructor) {
	self[key] = constructor
}

// Call a constructor for specified key or "*", if defined. Otherwise return an error
func (self Factory) Create(key string) (Parsable, error) {
	if constructor, ok := self[key]; ok {
		return constructor(), nil
	}

	// This is default constructor if defined
	if constructor, ok := self["*"]; ok {
		return constructor(), nil
	}
	return nil, errors.New("No element constructor defined for " + key)
}
