package item

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/extensions/iq/query"
	"github.com/dotdoom/goxmpp/stream/elements"
)

func init() {
	query.QueryFactory.AddConstructor("http://jabber.org/protocol/muc#admin item", func() elements.Element {
		return &Item{}
	})
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Affiliation string   `xml:"affiliation,attr"`
	JID         string   `xml:"jid,attr"`
	Text        string   `xml:",chardata"`
}
