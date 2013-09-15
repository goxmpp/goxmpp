package stream

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/decoder"
	"io"
	"log"
)

type ElementHandlerAction func(Element) bool

type InnerElementAdder interface {
	AddInnerElement(Element) bool
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

type InnerXMLHandler interface {
	InnerElementAdder
	HandleInnerXML(*Wrapper) []Element
}

type InnerXML struct {
	InnerElements  `xml:"omitempty"`
	InnerXML       []byte         `xml:",innerxml"`
	ElementFactory ElementFactory `xml:"-"`
}

func (self *InnerXML) Erase() {
	self.InnerXML = self.InnerXML[:0]
}

func (self *InnerXML) HandleInnerXML(sw *Wrapper) []Element {
	handlers := make([]Element, 0)

	if len(self.InnerXML) > 0 {
		sw.InnerDecoder.PutXML(self.InnerXML)

		processStreamElements(sw.InnerDecoder, self.ElementFactory, func(handler Element) bool {
			handlers = append(handlers, handler)
			return true
		})
	}
	self.Erase()

	return handlers
}

type XMLDecoder interface {
	Token() (xml.Token, error)
	DecodeElement(interface{}, *xml.StartElement) error
}

func processStreamElements(xmldecoder XMLDecoder, factory ElementFactory, elementAction ElementHandlerAction) {
	var token xml.Token
	var terr error

	for token, terr = xmldecoder.Token(); terr == nil; token, terr = xmldecoder.Token() {
		if element, ok := token.(xml.StartElement); ok && element.Name.Local != decoder.TERMINATOR {
			var handler Element
			var err error

			if handler, err = factory.Create(element.Name.Space + " " + element.Name.Local); err != nil {
				// TODO: added logging here
				continue
			}

			if err = xmldecoder.DecodeElement(handler, &element); err != nil {
				// TODO: added logging here
				log.Println(err)
				continue
			}

			if !elementAction(handler) {
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

func unmarshalStreamElement(self Element, sw *Wrapper) Element {
	// For elements other than InnerXMLHandler consider they don't have InnerElements
	if adder, ok := self.(InnerXMLHandler); ok {
		for _, element := range adder.HandleInnerXML(sw) {
			adder.AddInnerElement(unmarshalStreamElement(element, sw))
		}
	}
	return self
}
