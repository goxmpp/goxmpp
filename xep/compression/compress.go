package compression

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"

type Compression struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl compression"`
	stream.InnerElements
}
