package decoder

import "github.com/dotdoom/goxmpp/stream"

type ElementHandler interface {
	HandleElement(stream.Wrapper)
}

type InnerXMLHandler interface {
	HandleInnerXML(stream.Wrapper) []ElementHandler
}

type InnerXML struct {
	InnerXML []byte `xml:",innerxml"`
	extensions.Registrator
}

func (self *InnerXML) HandleInnerXML(sw stream.Wrapper) []ElementHandler {
	sw.InnerDecoder.PutXML(self.InnerXML)

	handlers := make([]ElementHandler)
	for token, terr := sw.InnerDecoder.Token(); err == nil; token, terr := sw.InnerDecoder.Token() {
		switch element, realType := token.(type); realType {
		case xml.StartElement:
			if handler, err := self.Registrator.GetHandler(elemnt.Name.Space + " " + element.Name.Local); err == nil {
				handlers = append(handlers, handler)
			}
		}
	}

	return handlers
}

func (self *InnerXML) HandleElement(sw stream.Wrapper) {
	for _, element := range self.HandlerInnerXML(sw) {
		element.HandleElement(sw)
	}
}
