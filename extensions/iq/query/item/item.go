package item

import (
	"encoding/xml"

	"github.com/goxmpp/xtream"
)

func init() {
	xtream.NodeFactory.Add(func() xtream.Element {
		return &Item{}
	}, xml.Name{Local: "query"}, xml.Name{Local: "item"})
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Affiliation string   `xml:"affiliation,attr"`
	JID         string   `xml:"jid,attr"`
	Text        string   `xml:",chardata"`
}
