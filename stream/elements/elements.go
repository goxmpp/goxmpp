package elements

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/decoder"
	"io"
	"log"
)

type Element interface{}

// Callback function to be called on each sibling element.
// Returns false to stop further processing.
type ElementParsedCallback func(Element) bool

type ElementAdder interface {
	AddElement(Element) bool
}

// Containers implementing this interface may accept inner elements,
// parsed by ParseElements.
type ParsableElementsContainer interface {
	ElementAdder
	ParseElements(*decoder.InnerDecoder) []Element
}

type Elements struct {
	Elements []Element
}

func (self *Elements) AddElement(e Element) bool {
	if e != nil {
		self.Elements = append(self.Elements, e)
		return true
	}
	return false
}

type UnmarshallableElements struct {
	Elements
	InnerXML       []byte  `xml:",innerxml"`
	ElementFactory Factory `xml:"-"`
}

func (self *UnmarshallableElements) ParseElements(decoder *decoder.InnerDecoder) (elements []Element) {
	if len(self.InnerXML) > 0 {
		decoder.PutXML(self.InnerXML)

		UnmarshalSiblingElements(decoder, self.ElementFactory, func(element Element) bool {
			elements = append(elements, element)
			return true
		})
	}
	// Reset InnerXML to avoid duplicates when marshalling
	self.InnerXML = nil
	return
}

// Recursively parse an element.
func ParseElement(self Element, decoder *decoder.InnerDecoder) Element {
	if adder, ok := self.(ParsableElementsContainer); ok {
		for _, element := range adder.ParseElements(decoder) {
			adder.AddElement(ParseElement(element, decoder))
		}
	}
	return self
}

type XMLDecoder interface {
	Token() (xml.Token, error)
	DecodeElement(interface{}, *xml.StartElement) error
}

// Run through the linear loop:
//  - get next sibling tag name from xmldecoder
//  - call factory to create an appropriate empty element instance
//  - unmarshal XML element into language object created by factory using DecodeElement
//  - pass parsed object to callback, breaking the loop if it returns false
func UnmarshalSiblingElements(xmldecoder XMLDecoder, factory Factory, callback ElementParsedCallback) {
	var token xml.Token
	var terr error

	for token, terr = xmldecoder.Token(); terr == nil; token, terr = xmldecoder.Token() {
		if xml_element, ok := token.(xml.StartElement); ok && xml_element.Name.Local != decoder.TERMINATOR {
			var obj_element Element
			var err error

			if obj_element, err = factory.Create(xml_element.Name.Space + " " + xml_element.Name.Local); err != nil {
				log.Println("Factory:", err, "(skipping)")
				continue
			}

			if err = xmldecoder.DecodeElement(obj_element, &xml_element); err != nil {
				// This is fatal. Return.
				log.Println("Decode:", err, "(skipping)")
				continue
			}

			if !callback(obj_element) {
				break
			}
		}

		if innerDecoder, ok := xmldecoder.(*decoder.InnerDecoder); ok && innerDecoder.IsEmpty() {
			break
		}
	}

	if terr != nil && terr != io.EOF {
		// This is fatal.
		log.Println(terr)
	}
}
