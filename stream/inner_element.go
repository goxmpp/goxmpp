package stream

import (
	"encoding/xml"
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
	InnerElements `xml:"omitempty"`
	InnerXML      []byte                    `xml:",innerxml"`
	Registrator   ElementHandlerRegistrator `xml:"-"`
}

func (self *InnerXML) Erase() {
	self.InnerXML = self.InnerXML[:0]
}

func (self *InnerXML) HandleInnerXML(sw *Wrapper) []Element {
	handlers := make([]Element, 0)

	if len(self.InnerXML) > 0 {
		sw.InnerDecoder.PutXML(self.InnerXML)

		processStreamElements(sw.InnerDecoder.Decoder, self.Registrator, func(handler Element) bool {
			handlers = append(handlers, handler)
			return len(*sw.InnerDecoder.InnerXMLBuffer) > 0
		})
	}
	self.Erase()

	return handlers
}

func processStreamElements(decoder *xml.Decoder, registry ElementHandlerRegistrator, elementAction ElementHandlerAction) {
	var token xml.Token
	var terr error

OUT:
	for token, terr = decoder.Token(); terr == nil; token, terr = decoder.Token() {
		switch element := token.(type) {
		case xml.StartElement:
			var handler Element
			var err error

			if handler, err = registry.GetHandler(element.Name.Space + " " + element.Name.Local); err != nil {
				// TODO: added logging here
				continue
			}

			if err = decoder.DecodeElement(handler, &element); err != nil {
				// TODO: added logging here
				continue
			}

			if !elementAction(handler) {
				break OUT
			}
		}
	}

	if terr != nil {
		// TODO: log error
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
