package elements

import (
	"encoding/xml"
	"log"
)

// Create an (empty) Parsable to parse XML into
type Constructor func() Element

// Maintain a mapping between tag names (and namespaces) and Constructors
type Factory map[string]Constructor

func NewFactory() Factory {
	return Factory(make(map[string]Constructor))
}

func (self Factory) AddConstructor(key string, constructor Constructor) {
	log.Printf("Factory: adding constructor for %#v\n", key)
	self[key] = constructor
}

// Call a constructor for specified key or "*", if defined. Otherwise return an error
func (self Factory) Get(element xml.StartElement) (Element, error) {
	name_key := " " + element.Name.Local
	full_key := element.Name.Space + name_key

	log.Printf("Factory: searching for fully-qualified key %#v\n", full_key)
	if constructor, ok := self[full_key]; ok {
		return constructor(), nil
	}

	log.Printf("Factory: full key missed, searching for local key %#v\n", name_key)
	if constructor, ok := self[name_key]; ok {
		return constructor(), nil
	}

	// This is default constructor if defined
	log.Printf("Factory: local key missed, attempting default\n")
	if constructor, ok := self["*"]; ok {
		return constructor(), nil
	}

	return &InnerXML{}, nil
}

func (self Factory) DecodeElement(d *xml.Decoder, element *xml.StartElement) (interface{}, error) {
	elementObject, err := self.Get(*element)
	if err != nil {
		return nil, err
	}
	if err := d.DecodeElement(elementObject, element); err != nil {
		return nil, err
	}

	return elementObject, nil
}
