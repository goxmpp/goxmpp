package query

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/stanzas/iq"
)

func init() {
	iq.Factory.AddConstructor("http://jabber.org/protocol/muc#admin query", func() elements.Element {
		return &MucQuery{InnerElements: elements.NewInnerElements(Factory)}
	})
	iq.Factory.AddConstructor("http://jabber.org/protocol/disco#info query", func() elements.Element {
		return &DiscoQuery{InnerElements: elements.NewInnerElements(Factory)}
	})
}

var Factory = elements.NewFactory()

type MucQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#admin query"`
	*elements.InnerElements
}

type DiscoQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#info query"`
	*elements.InnerElements
}
