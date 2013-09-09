package stream

import "encoding/xml"
import "github.com/dotdoom/goxmpp"
import . "github.com/dotdoom/goxmpp/interfaces"

type Stream struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams stream"`
	ID      string   `xml:"id,attr,omitempty"`
	From    string   `xml:"from,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	Version string   `xml:"version,attr,omitempty"`
}

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
}

func (self *InnerXML) HandleInnerXML(sw goxmpp.StreamWrapper) []ElementHandler {
	// TODO: Put innerXML parser here
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
