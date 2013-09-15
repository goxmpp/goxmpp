package query

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/iq"
)

func init() {
	iq.ElementFactory.AddConstructor("http://jabber.org/protocol/muc#admin query", func() elements.Element {
		return &MucQuery{InnerXML: stream.InnerXML{ElementFactory: ElementFactory}}
	})
	iq.ElementFactory.AddConstructor("http://jabber.org/protocol/disco#info query", func() elements.Element {
		return &DiscoQuery{InnerXML: stream.InnerXML{ElementFactory: ElementFactory}}
	})
}

var ElementFactory = elements.NewFactory()

type MucQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#admin query"`
	stream.InnerXML
}

type DiscoQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#info query"`
	stream.InnerXML
}
