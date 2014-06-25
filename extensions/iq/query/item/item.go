package item

import (
	"encoding/xml"

	"github.com/goxmpp/goxmpp/extensions/iq/query"
	"github.com/goxmpp/goxmpp/stream/elements"
)

func init() {
	query.QueryFactory.AddConstructor(func() elements.Element {
		return &Item{}
	})
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Affiliation string   `xml:"affiliation,attr"`
	JID         string   `xml:"jid,attr"`
	Text        string   `xml:",chardata"`
}
