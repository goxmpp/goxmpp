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

type InnerElementsAdder interface {
	AddInnerElement(Element) bool
}

// Containers implementing this interface may accept inner elements,
// parsed by unmarshalStreamElement
type ParsedInnerElementsContainer interface {
	InnerElementsAdder
	ParseInnerElements(*decoder.InnerDecoder) []Element
}

type InnerElements struct {
	InnerElements []Element
}

func (self *InnerElements) AddInnerElement(e Element) bool {
	if e != nil {
		self.InnerElements = append(self.InnerElements, e)
		return true
	}
	return false
}

type InnerXML struct {
	InnerElements
	InnerXML       []byte  `xml:",innerxml"`
	ElementFactory Factory `xml:"-"`
}

func (self *InnerXML) ParseInnerElements(decoder *decoder.InnerDecoder) (elements []Element) {
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

// Recursively unmarshal an element.
func UnmarshalElement(self Element, decoder *decoder.InnerDecoder) Element {
	// For elements other than InnerXMLParser consider they don't have InnerElements
	if adder, ok := self.(ParsedInnerElementsContainer); ok {
		for _, element := range adder.ParseInnerElements(decoder) {
			adder.AddInnerElement(UnmarshalElement(element, decoder))
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
				log.Println(err)
				continue
			}

			if err = xmldecoder.DecodeElement(obj_element, &xml_element); err != nil {
				log.Println(err)
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
		log.Println(terr)
	}
}
