package stream

type ElementHandler interface {
	HandleElement(Wrapper)
}

type InnerXMLHandler interface {
	HandleInnerXML(Wrapper) []ElementHandler
}

type InnerXML struct {
	InnerXML []byte `xml:",innerxml"`
	extensions.Registrator
}

func (self *InnerXML) HandleInnerXML(sw Wrapper) []ElementHandler {
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

func (self *InnerXML) HandleElement(sw Wrapper) {
	for _, element := range self.HandlerInnerXML(sw) {
		element.HandleElement(sw)
	}
}
