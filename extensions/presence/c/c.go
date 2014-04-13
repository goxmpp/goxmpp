package c

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/stanzas/presence"
)

type CElement struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/caps c"`
	Node    string   `xml:"node,attr"`
	Ver     string   `xml:"ver,attr"`
	Hash    string   `xml:"hash,attr"`
	Ext     string   `xml:"ext,attr,omitempty"`
}

func init() {
	presence.PresenceFactory.AddConstructor(func() elements.Element {
		return &CElement{}
	})
}
