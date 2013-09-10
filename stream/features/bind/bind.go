package bind

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"

type Feature struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	stream.InnerElements
}
