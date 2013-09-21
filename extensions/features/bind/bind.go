package bind

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type Feature struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	features.Elements
}
