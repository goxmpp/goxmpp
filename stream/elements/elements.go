package elements

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/decoder"
	"io"
	"log"
)

// An interface containing enough methods for parsing/unmarshalling to work.
type XMLDecoder interface {
	Token() (xml.Token, error)
	DecodeElement(interface{}, *xml.StartElement) error
}

// Run through the linear loop:
//  - get next sibling tag name (and xmlns) from xmldecoder.Token()
//  - call factory to create an empty element instance for that tag name
//  - use Parse() on the created empty element
//  - pass parsed object to callback, breaking the loop if it returns false
func ParseSiblingElements(xmldecoder XMLDecoder, factory Factory, callback ElementParsedCallback) error {
	var token xml.Token
	var fatal error

	for token, fatal = xmldecoder.Token(); fatal == nil; token, fatal = xmldecoder.Token() {
		if xml_element, ok := token.(xml.StartElement); ok && xml_element.Name.Local != decoder.TERMINATOR {
			var obj_element Parsable
			var err error

			if obj_element, err = factory.Create(xml_element.Name.Space + " " + xml_element.Name.Local); err != nil {
				log.Println("Factory:", err, "(skipping)")
				continue
			}

			if err = obj_element.Parse(xmldecoder, &xml_element); err != nil {
				fatal = err
				break
			}

			if !callback(obj_element) {
				break
			}
		}

		if innerDecoder, ok := xmldecoder.(*decoder.InnerDecoder); ok && innerDecoder.IsEmpty() {
			break
		}
	}

	if fatal == io.EOF {
		return nil
	}
	return fatal
}

// Callback function to be called on each sibling element unmarshalled.
// This function is expected to continue parsing (not necessarily unmarshalling).
// Returns false to stop further processing.
type ElementParsedCallback func(Parsable) bool

// structs implementing this interface can parse themselves.
// They are typically registered within a factory.
type Parsable interface {
	Parse(XMLDecoder, *xml.StartElement) error
}

// Parsable implementation using unmarshalling (xml.DecodeElement)
type Unmarshallable struct {
}

func (self *Unmarshallable) Parse(decoder XMLDecoder, start *xml.StartElement) error {
	return decoder.DecodeElement(self, start)
}

// structs implementing this interface can parse inner elements.
type ParsableElements interface {
	ParseElements(*decoder.InnerDecoder) error
}

// ParsableElements implementation based on parsing contents of a
// previously unmarshalled field tagged as 'innerxml'. Elements are stored in a slice.
type ElementsFromInnerXML struct {
	Elements       []Parsable
	InnerXML       []byte  `xml:",innerxml"`
	ElementFactory Factory `xml:"-"`
}

// Unmarshall all elements in InnerXML, then call ParseElements on each individually.
// Can't go for immediate recursion because it will replace the XML stored in decoder.
func (self *ElementsFromInnerXML) ParseElements(decoder *decoder.InnerDecoder) error {
	var fatal error

	if len(self.InnerXML) > 0 {
		decoder.PutXML(self.InnerXML)

		fatal = ParseSiblingElements(decoder, self.ElementFactory, func(element Parsable) bool {
			self.Elements = append(self.Elements, element)
			return true
		})

		decoder.PutXML(nil)

		if fatal != nil {
			return fatal
		}
	}

	// Reset InnerXML to avoid duplicates when marshalling
	self.InnerXML = nil

	for _, element := range self.Elements {
		if element, ok := element.(ParsableElements); ok {
			if fatal = element.ParseElements(decoder); fatal != nil {
				return fatal
			}
		}
	}

	return nil
}

// Simple encoder interface
type XMLEncoder interface {
	Encode(interface{}) error
}

// structs implementing this interface can compose themselves into
// xml.Encoder or io.Writer
type Composable interface {
	Compose(XMLEncoder, io.Writer) error
}

// Composable implementation using marshalling.
type Marshallable struct {
}

func (self *Marshallable) Compose(encoder XMLEncoder, _ io.Writer) error {
	return encoder.Encode(self)
}
