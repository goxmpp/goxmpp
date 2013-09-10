package privacy

import "encoding/xml"

type PrivacyQuery struct {
	// http://xmpp.org/rfcs/rfc3921.html
	XMLName xml.Name `xml:"jabber:iq:privacy query"`
}
