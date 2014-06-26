package query

import (
	"encoding/xml"

	"github.com/goxmpp/xtream"
)

func init() {
	iqXMLName := xml.Name{Local: "iq"}
	discoXMLName := xml.Name{Local: "query", Space: "http://jabber.org/protocol/disco#info"}
	mucXMLName := xml.Name{Local: "query", Space: "http://jabber.org/protocol/muc#admin"}

	xtream.NodeFactory.Add(func() xtream.Element {
		return &MucQuery{InnerElements: xtream.NewElements(&mucXMLName)}
	}, iqXMLName, mucXMLName)
	xtream.NodeFactory.Add(func() xtream.Element {
		return &DiscoQuery{InnerElements: xtream.NewElements(&discoXMLName)}
	}, iqXMLName, discoXMLName)
}

type MucQuery struct {
	XMLName              xml.Name `xml:"http://jabber.org/protocol/muc#admin query"`
	xtream.InnerElements `xml:",any"`
}

type DiscoQuery struct {
	XMLName              xml.Name `xml:"http://jabber.org/protocol/disco#info query"`
	xtream.InnerElements `xml:",any"`
}
