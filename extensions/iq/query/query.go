package query

import (
	"encoding/xml"

	"github.com/goxmpp/xtream"
)

func init() {
	xtream.NodeFactory.Add(func(name *xml.Name) xtream.Element {
		return &MucQuery{InnerElements: xtream.NewElements(name)}
	})
	xtream.NodeFactory.Add(func(name *xml.Name) xtream.Element {
		return &DiscoQuery{InnerElements: xtream.NewElements(name)}
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
