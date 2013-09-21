package starttls

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type StartTLSStreamFeature struct {
	XMLName     xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required    bool     `xml:"required,omitempty"`
	Certificate []byte
	features.Elements
}
