package auth

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"

type Mechanisms struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	stream.InnerElements
}

type AuthMechanismStreamFeature struct {
	XMLName xml.Name `xml:"mechanism"`
	Name    string   `xml:",chardata"`
	stream.InnerElements
}
