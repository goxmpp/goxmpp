package stanza

import "github.com/dotdoom/goxmpp"

type InnerElementAdder interface {
	AddInnerElement(InnerElementAdder) bool
}

type ElementHandler interface {
	HandleElement(goxmpp.StreamWrapper)
}

type XMLHandler interface {
	HandleInnerXML(goxmpp.StreamWrapper) []ElementHandler
}

type InnerElements struct {
	InnerElements []interface{}
}

func (self *InnerElements) AddInnerElement(e interface{}) bool {
	if e != nil {
		self.InnerElements = append(self.InnerElements, e)
		return true
	}
	return false
}

type Stanza struct {
	From string `xml:"from,attr,omitempty"`
	To   string `xml:"to,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	ID   string `xml:"id,attr,omitempty"`
	Lang string `xml:"xml:lang,attr,omitempty"`
	InnerElements
}

type InnerXML struct {
	XML []byte `xml:"innerxml"`
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

	for _, handler := range handlers {
		handler.Handle()
	}

	return make([]ElementHandler)
}

func (self *InnerXML) HandleElement(sw goxmpp.StreamWrapper) {
	for _, element := range self.HandlerInnerXML(sw) {
		element.HandleElement(sw)
	}
}
