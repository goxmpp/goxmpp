package query

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/iq"
)

func init() {
	iq.HandlerRegistrator.Register("http://jabber.org/protocol/muc#admin query", func() stream.Element {
		return &MucQuery{InnerXML: stream.InnerXML{Registrator: HandlerRegistrator}}
	})
	iq.HandlerRegistrator.Register("http://jabber.org/protocol/disco#info query", func() stream.Element {
		return &DiscoQuery{InnerXML: stream.InnerXML{Registrator: HandlerRegistrator}}
	})
}

var HandlerRegistrator = stream.NewElementGeneratorRegistrator()

type MucQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#admin query"`
	stream.InnerXML
}

type DiscoQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#info query"`
	stream.InnerXML
}
