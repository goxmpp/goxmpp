package bind

import "encoding/xml"

type Feature struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	stream.InnerElements
}
