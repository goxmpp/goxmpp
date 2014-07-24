package time

import "encoding/xml"
import "github.com/goxmpp/xtream"

func init() {
	xtream.NodeFactory.Add(func(name *xml.Name) xtream.Element { return &TimeQuery{} })
}

type TimeQuery struct {
	// http://xmpp.org/extensions/xep-0202.html
	XMLName xml.Name `xml:"urn:xmpp:time time" parent:"iq"`
	TZO     string   `xml:"tzo,omitempty"`
	UTC     string   `xml:"utc,omitempty"`
}
