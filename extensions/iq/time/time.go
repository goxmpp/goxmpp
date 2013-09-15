package time

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream/elements"
import "github.com/dotdoom/goxmpp/stream/elements/iq"

func init() {
	iq.ElementFactory.AddConstructor("urn:xmpp:time time", func() elements.Element { return &TimeQuery{} })
}

type TimeQuery struct {
	// http://xmpp.org/extensions/xep-0202.html
	XMLName xml.Name `xml:"urn:xmpp:time time"`
	TZO     string   `xml:"tzo,omitempty"`
	UTC     string   `xml:"utc,omitempty"`
}
