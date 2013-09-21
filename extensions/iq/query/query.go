package query

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/stanzas/iq"
)

func init() {
	iq.ElementFactory.AddConstructor("http://jabber.org/protocol/muc#admin query", func() elements.Element {
		return &MucQuery{UnmarshallableElements: elements.UnmarshallableElements{ElementFactory: ElementFactory}}
	})
	iq.ElementFactory.AddConstructor("http://jabber.org/protocol/disco#info query", func() elements.Element {
		return &DiscoQuery{UnmarshallableElements: elements.UnmarshallableElements{ElementFactory: ElementFactory}}
	})
}

var ElementFactory = elements.NewFactory()

type MucQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#admin query"`
	elements.UnmarshallableElements
}

type DiscoQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#info query"`
	elements.UnmarshallableElements
}
