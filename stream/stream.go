package stream

import "encoding/xml"

type Stream struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams stream"`
	ID      string   `xml:"id,attr,omitempty"`
	From    string   `xml:"from,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	Version string   `xml:"version,attr,omitempty"`
}

type InnerElementAdder interface {
	AddSubElement(InnerElementAdder) bool
}

type InnerElements struct {
	Elements []InnerElementAdder
}

func (self *InnerElements) AddSubElement(se InnerElementAdder) bool {
	if sf != nil {
		self.Elements = append(self.Elements, se)
		return true
	}
	return false
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
