package stanza

import "github.com/dotdoom/goxmpp"

type InnerElementAdder interface {
	AddInnerElement(InnerElementAdder) bool
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

