package stream

import "encoding/xml"
import "github.com/dotdoom/goxmpp"
import "github.com/dotdoom/goxmpp/xep"
import . "github.com/dotdoom/goxmpp/interfaces"

type Stream struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams stream"`
	ID      string   `xml:"id,attr,omitempty"`
	From    string   `xml:"from,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	Version string   `xml:"version,attr,omitempty"`
}

var Registrator = xep.NewRegistrator()

type InnerElements struct {
	Elements []InnerElementAdder
}

func (self *InnerElements) AddInnerElement(se InnerElementAdder) bool {
	if sf != nil {
		self.Elements = append(self.Elements, se)
		return true
	}
	return false
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

type Stanza struct {
	From  string `xml:"from,attr,omitempty"`
	To    string `xml:"to,attr,omitempty"`
	Type  string `xml:"type,attr,omitempty"`
	ID    string `xml:"id,attr,omitempty"`
	Error Error
	InnerElements
}

type Error struct {
	Type string `xml:"type,attr,omitempty"`
}
