package ping

import "encoding/xml"

type PingQuery struct {
	// http://xmpp.org/extensions/xep-0199.html
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}
