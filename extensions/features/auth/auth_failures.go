package auth

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/stream/elements"
)

type Failure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl failure"`
	elements.Element
}

func NewFailute(el elements.Element) Failure {
	return Failure{Element: el}
}

type NotAuthorized struct {
	XMLName xml.Name `xml:"not-authorized"`
}

type InvalidMechanism struct {
	XMLName xml.Name `xml:"invalid-mechanism"`
}
