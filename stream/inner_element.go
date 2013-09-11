package stream

import (
	"encoding/xml"
)

type ElementHandlerAction func(ElementHandler)

type ElementHandler interface {
	HandleElement(*Wrapper)
}

type InnerXMLHandler interface {
	HandleInnerXML(*Wrapper) []ElementHandler
}

type InnerXML struct {
	InnerXML    []byte `xml:",innerxml"`
	Registrator ElementHandlerRegistrator
}

func (self *InnerXML) HandleInnerXML(sw *Wrapper) []ElementHandler {
	sw.InnerDecoder.PutXML(self.InnerXML)
	handlers := make([]ElementHandler, 0)

	processStreamElements(sw.InnerDecoder.Decoder, self.Registrator, func(handler ElementHandler) {
		handlers = append(handlers, handler)
	})

	return handlers
}

func processStreamElements(decoder *xml.Decoder, registry ElementHandlerRegistrator, elementAction ElementHandlerAction) {
	for token, terr := decoder.Token(); terr == nil; token, terr = decoder.Token() {
		switch element := token.(type) {
		case xml.StartElement:
			var handler ElementHandler
			var err error
			if handler, err = registry.GetHandler(element.Name.Space + " " + element.Name.Local); err != nil {
				// TODO: added logging here
				continue
			}

			if err = decoder.DecodeElement(handler, &element); err != nil {
				// TODO: added logging here
				continue
			}

			elementAction(handler)
		}
	}
}

func (self *InnerXML) HandleElement(sw *Wrapper) {
	for _, element := range self.HandleInnerXML(sw) {
		element.HandleElement(sw)
	}
}
