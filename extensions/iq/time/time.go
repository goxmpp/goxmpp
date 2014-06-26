package time

import "encoding/xml"
import "github.com/goxmpp/xtream"

func init() {
	xtream.NodeFactory.Add(func() xtream.Element { return &TimeQuery{} }, xml.Name{Local: "iq"}, xml.Name{Local: "time"})
}

type TimeQuery struct {
	// http://xmpp.org/extensions/xep-0202.html
	XMLName xml.Name `xml:"urn:xmpp:time time"`
	TZO     string   `xml:"tzo,omitempty"`
	UTC     string   `xml:"utc,omitempty"`
}
