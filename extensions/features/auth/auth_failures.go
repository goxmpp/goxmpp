package auth

import (
	"encoding/xml"

	"github.com/goxmpp/xtream"
)

type Failure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl failure"`
	xtream.Element
}

func NewFailute(el xtream.Element) Failure {
	return Failure{Element: el}
}

type NotAuthorized struct {
	XMLName xml.Name `xml:"not-authorized"`
}

type InvalidMechanism struct {
	XMLName xml.Name `xml:"invalid-mechanism"`
}
