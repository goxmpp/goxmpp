package last

import "encoding/xml"

type LastQuery struct {
	// http://xmpp.org/extensions/xep-0012.html
	XMLName xml.Name `xml:"jabber:iq:last query"`
	Seconds int      `xml:"seconds,attr,omitempty"`
}
