package stanzas

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/stream/elements"
)

type Base struct {
	From string `xml:"from,attr,omitempty"`
	To   string `xml:"to,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	ID   string `xml:"id,attr,omitempty"`
	Lang string `xml:"xml:lang,attr,omitempty"`
}

func (b *Base) SetFromStartElement(start xml.StartElement) {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "from":
			b.From = attr.Value
		case "to":
			b.To = attr.Value
		case "type":
			b.Type = attr.Value
		case "id":
			b.ID = attr.Value
		case "lang":
			b.Lang = attr.Value
		}
	}
}

var Factory = elements.NewElementFactory()
