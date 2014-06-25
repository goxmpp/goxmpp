package query

import (
	"encoding/xml"

	"github.com/goxmpp/goxmpp/stream/elements"
	"github.com/goxmpp/goxmpp/stream/elements/stanzas/iq"
)

func init() {
	iq.IQFactory.AddConstructor(func() elements.Element {
		return &MucQuery{InnerElements: elements.NewInnerElements(QueryFactory)}
	})
	iq.IQFactory.AddConstructor(func() elements.Element {
		return &DiscoQuery{InnerElements: elements.NewInnerElements(QueryFactory)}
	})
}

var QueryFactory = elements.NewFactory()

type MucQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#admin query"`
	*elements.InnerElements
}

type DiscoQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#info query"`
	*elements.InnerElements
}
