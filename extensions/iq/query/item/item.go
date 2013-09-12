package item

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/extensions/iq/query"
	"github.com/dotdoom/goxmpp/stream"
)

func init() {
	query.HandlerRegistrator.Register(" item", func() stream.Element {
		return &Item{}
	})
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Affiliation string   `xml:"affiliation,attr"`
	JID         string   `xml:"jid,attr"`
	Text        string   `xml:",chardata"`
}
