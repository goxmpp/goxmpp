package starttls

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

func init() {
	features.Tree.AddElement(NewStartTLSFeature())
	features.Factory.AddConstructor("urn:ietf:params:xml:ns:xmpp-tls starttls", func() elements.Element {
		return &StartTLSElement{}
	})
}

type StartTLSFeatureElement struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required bool     `xml:"required,omitempty"`
}

func NewStartTLSFeature() *StartTLSFeatureElement {
	return &StartTLSFeatureElement{}
}

type StartTLSElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
}

func (s *StartTLSElement) Handle(stream *stream.Stream) error {
	stream.WriteElement(&ProceedElement{})

}

type ProceedElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}
