package c

import (
	"encoding/xml"

	"github.com/goxmpp/xtream"
)

type CElement struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/caps c" parent:"presence"`
	Node    string   `xml:"node,attr"`
	Ver     string   `xml:"ver,attr"`
	Hash    string   `xml:"hash,attr"`
	Ext     string   `xml:"ext,attr,omitempty"`
}

func init() {
	xtream.NodeFactory.Add(func() xtream.Element {
		return &CElement{}
	})
}
