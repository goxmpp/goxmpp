package starttls

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"

type StartTLSStreamFeature struct {
	XMLName     xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required    bool     `xml:"required,omitempty"`
	Certificate []byte
	stream.InnerElements
}
