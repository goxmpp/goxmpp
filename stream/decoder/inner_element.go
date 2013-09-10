package decoder

import "github.com/dotdoom/goxmpp/stream"

type ElementHandler interface {
	HandleElement(stream.Wrapper)
}

type InnerXMLHandler interface {
	HandleInnerXML(stream.Wrapper) []ElementHandler
}

type InnerXML struct {
	XML []byte `xml:",innerxml"`
	xep.Registrator
}

func (self *InnerXML) HandleInnerXML(sw goxmpp.StreamWrapper) []ElementHandler {
	sw.Buffer.PutXML(self.XML)

	handlers := make([]ElementHandler)
	for token, terr := sw.Decoder.Token(); err == nil; token, terr := sw.Decoder.Token() {
		switch element, realType := token.(type); realType {
		case xml.StartElement:
			if handler, err := self.Registrator.GetHandler(elemnt.Name.Space + " " + element.Name.Local); err == nil {
				handlers = append(handlers, handler)
			}
		}
	}

	return handlers
}

func (self *InnerXML) HandleElement(sw goxmpp.StreamWrapper) {
	for _, element := range self.HandlerInnerXML(sw) {
		element.HandleElement(sw)
	}
}
