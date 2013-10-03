package auth

import (
	"encoding/xml"
	_ "github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type mechanismsElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	features.Elements
}

func (self *mechanismsElement) IsRequiredFor(fs features.State) bool {
	return fs["authenticated"] == nil
}

func (self *mechanismsElement) CopyIfAvailable(fs features.State) interface{} {
	if self.IsRequiredFor(fs) {
		return self.CopyAvailableFeatures(fs, new(mechanismsElement))
	}
	return nil
}

type MechanismElement struct {
	XMLName xml.Name `xml:"mechanism"`
	Name    string   `xml:",chardata"`
	features.Elements
}

var Mechanisms = new(mechanismsElement)

type AuthElement struct {
	XMLName   xml.Name `xml:"auth"`
	Mechanism string   `xml:"mechanism,attr"`
	Data      string   `xml:",chardata"`
	elements.UnmarshallableElements
}

type AuthElementHandler interface {
	Handle(*AuthElement, features.State, interface{}) bool
}

func (self *AuthElement) HandleFeature(state features.State, sw interface{}) {
	for _, m := range Mechanisms.Elements.Elements {
		if m.(AuthElementHandler).Handle(self, state, sw) {
			break
		}
	}
}

var ElementFactory = elements.NewFactory()

func init() {
	features.List.AddElement(Mechanisms)
	features.Factory.AddConstructor("urn:ietf:params:xml:ns:xmpp-sasl auth", func() elements.Element {
		return &AuthElement{UnmarshallableElements: elements.UnmarshallableElements{ElementFactory: ElementFactory}}
	})
}
