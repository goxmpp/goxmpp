package auth

import (
	"encoding/xml"
	"fmt"
	"log"

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

func NewAuthElement() *AuthElement {
	return &AuthElement{}
}

func (self *AuthElement) Handle(stream *stream.Stream) error {
	log.Printf("handling auth %s\n", self.Mechanism)
	if handler := Mechanisms[self.Mechanism]; handler != nil {
		return handler(self, stream)
	} else {
		return fmt.Errorf("No handler for mechanism %v", self.Mechanism)
	}
}

var Factory = elements.NewElementFactory()

var Mechanisms map[string]Handler = make(map[string]Handler)

func init() {
	features.Factory.AddConstructor("urn:ietf:params:xml:ns:xmpp-sasl auth", func() elements.Element {
		log.Println("creating new AuthElement")
		return NewAuthElement()
	})
	features.Tree.AddElement(Features)
}
