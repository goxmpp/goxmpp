package stanzas

import (
	_ "encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements"
)

type Base struct {
	From string `xml:"from,attr,omitempty"`
	To   string `xml:"to,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	ID   string `xml:"id,attr,omitempty"`
	Lang string `xml:"xml:lang,attr,omitempty"`
}

var Factory = elements.NewFactory()
