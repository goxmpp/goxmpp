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

var Mechanisms = new(mechanismsElement)

type MechanismElement struct {
	XMLName xml.Name `xml:"mechanism"`
	Name    string   `xml:",chardata"`
	features.Elements
}

type AuthElement struct {
	XMLName   xml.Name `xml:"auth"`
	Mechanism string   `xml:"mechanism,attr"`
	Data      string   `xml:",chardata"`
	elements.UnmarshallableElements
}

type Mechanism interface {
	Process(*AuthElement)
}

func (self *AuthElement) React(state features.State, conn features.SuperInterface) {
	println("Reacting on Auth, mechanism:", self.Mechanism, ", data:", self.Data)
	for _, m := range Mechanisms.Elements.Elements.Elements {
		m.(Mechanism).Process(self)
	}
	conn.NextElement()
}

var ElementFactory = elements.NewFactory()

func init() {
	features.List.AddElement(Mechanisms)
	features.Factory.AddConstructor("urn:ietf:params:xml:ns:xmpp-sasl auth", func() elements.Element {
		return &AuthElement{UnmarshallableElements: elements.UnmarshallableElements{ElementFactory: ElementFactory}}
	})
}
