package auth

import (
	"encoding/xml"
	"errors"
	"log"

	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type AuthElement struct {
	XMLName   xml.Name `xml:"auth"`
	Mechanism string   `xml:"mechanism,attr"`
	*elements.InnerElements
}

//
//type AuthElementHandler interface {
//	Handle(*AuthElement, features.List) bool
//}
//
//func (self *AuthElement) HandleFeature(state features.State, sw interface{}) {
//	//for _, m := range Mechanisms.Elements.Elements {
//	//	if m.(AuthElementHandler).Handle(self, state, sw) {
//	//		break
//	//	}
//	//}
//}

func NewAuthElement() *AuthElement {
	return &AuthElement{InnerElements: elements.NewInnerElements(Factory)}
}

func (self *AuthElement) Handle(stream *stream.Stream) error {
	log.Println("handling auth")
	if handler := Mechanisms[self.Mechanism]; handler != nil {
	}
	return errors.New("no handler found")
}

var Factory = elements.NewElementFactory()

var Mechanisms map[string]features.Handler

func init() {
	features.Factory.AddConstructor("urn:ietf:params:xml:ns:xmpp-sasl auth", func() elements.Element {
		return NewAuthElement()
	})
	features.Tree.AddElement(Features)
}
