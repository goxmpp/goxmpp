package auth

import (
	"encoding/xml"
	"fmt"

	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type AuthElement struct {
	XMLName   xml.Name `xml:"auth"`
	Mechanism string   `xml:"mechanism,attr"`
	Data      string   `xml:",chardata"`
}

type Handler func(*AuthElement, *stream.Stream) error

func (self *AuthElement) Handle(stream *stream.Stream) error {
	if handler := mechanism_handlers[self.Mechanism]; handler != nil {
		return handler(self, stream)
	} else {
		return fmt.Errorf("No handler for mechanism %v", self.Mechanism)
	}
}

func init() {
	features.Factory.AddConstructor("urn:ietf:params:xml:ns:xmpp-sasl auth", func() elements.Element {
		return &AuthElement{}
	})
}
