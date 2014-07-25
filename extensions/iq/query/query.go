package query

import (
	"encoding/xml"

	"github.com/goxmpp/xtream"
)

func init() {
	xtream.NodeFactory.Add(func() xtream.Element {
		return &MucQuery{InnerElements: xtream.NewElements()}
	})
	xtream.NodeFactory.Add(func() xtream.Element {
		return &DiscoQuery{InnerElements: xtream.NewElements()}
	})
}

type MucQuery struct {
	XMLName              xml.Name `xml:"http://jabber.org/protocol/muc#admin query" parent:"iq"`
	xtream.InnerElements `xml:",any"`
}

type DiscoQuery struct {
	XMLName              xml.Name `xml:"http://jabber.org/protocol/disco#info query" parent:"iq"`
	xtream.InnerElements `xml:",any"`
}
