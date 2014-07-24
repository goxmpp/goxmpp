package item

import (
	"encoding/xml"

	"github.com/goxmpp/xtream"
)

func init() {
	xtream.NodeFactory.Add(func(*xml.Name) xtream.Element {
		return &Item{}
	})
}

type Item struct {
	XMLName     xml.Name `xml:"item" parent:"query"`
	Affiliation string   `xml:"affiliation,attr"`
	JID         string   `xml:"jid,attr"`
	Text        string   `xml:",chardata"`
}
