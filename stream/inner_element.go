package stream

import (
	"encoding/xml"
)

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
	for token, terr := sw.InnerDecoder.Token(); terr == nil; token, terr = sw.InnerDecoder.Token() {
		switch element := token.(type) {
		case xml.StartElement:
			if handler, err := self.Registrator.GetHandler(element.Name.Space + " " + element.Name.Local); err == nil {
				handlers = append(handlers, handler)
			}
		}
	}

	return handlers
}

func (self *InnerXML) HandleElement(sw *Wrapper) {
	for _, element := range self.HandleInnerXML(sw) {
		element.HandleElement(sw)
	}
}
