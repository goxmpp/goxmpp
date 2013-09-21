package auth

import (
	"encoding/xml"
	_ "github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/connection"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type mechanisms struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	features.Elements
}

func (self *mechanisms) IsRequiredFor(fs connection.State) bool {
	return fs["authenticated"] == nil
}

func (self *mechanisms) CopyIfAvailable(fs connection.State) interface{} {
	if self.IsRequiredFor(fs) {
		return self.CopyAvailableFeatures(fs, new(mechanisms))
	}
	return nil
}

var Mechanisms = new(mechanisms)

type Mechanism struct {
	XMLName xml.Name `xml:"mechanism"`
	Name    string   `xml:",chardata"`
	features.Elements
}

type Auth struct {
	XMLName   xml.Name `xml:"auth"`
	Mechanism string   `xml:"mechanism,attr"`
	elements.UnmarshallableElements
}

func (self *Auth) React(conn *connection.Connection) {
	println("Reacting on: Auth")
}

var ElementFactory = elements.NewFactory()

func init() {
	features.List.AddElement(Mechanisms)
	features.Factory.AddConstructor("urn:ietf:params:xml:ns:xmpp-sasl auth", func() elements.Element {
		return &Auth{UnmarshallableElements: elements.UnmarshallableElements{ElementFactory: ElementFactory}}
	})
}
