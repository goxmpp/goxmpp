package time

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream/elements"
import "github.com/dotdoom/goxmpp/stream/elements/stanzas/iq"

func init() {
	iq.IQFactory.AddConstructor(func() elements.Element { return &TimeQuery{} })
}

type TimeQuery struct {
	// http://xmpp.org/extensions/xep-0202.html
	XMLName xml.Name `xml:"urn:xmpp:time time"`
	TZO     string   `xml:"tzo,omitempty"`
	UTC     string   `xml:"utc,omitempty"`
}
