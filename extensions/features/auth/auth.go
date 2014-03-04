package auth

import (
	"encoding/xml"
	_ "github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type AuthFeature interface {
	UserName() string
	SetUserName(username string)
	//features.StringElementMapper /*
	//	AddElement
	//*/
}

type mechanismsElement struct {
	XMLName    xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	username   string   `xml:"-"`
	mechanisms map[string]interface{}
	//features.StringElements
}

func (self *mechanismsElement) UserName() string {
	return self.username
}

func (self *mechanismsElement) SetUserName(username string) {
	self.username = username
}

func (self *mechanismsElement) IsRequired() bool {
	return self.username == ""
}

func (self *mechanismsElement) IsAvailable() bool {
	return self.IsRequired()
}

//func (self *mechanismsElement) AddMechanism(name string, mechanism interface{}) {
//	if self.mechanisms == nil {
//		self.mechanisms = make(map[string]interface{})
//	}
//	self.mechanisms[name] = mechanism
//}
//
//func (self *mechanismsElement) FindMechanism(name string) interface{} {
//	return self.mechanisms[name]
//}
//
//type MechanismElement struct {
//	XMLName xml.Name `xml:"mechanism"`
//	Name    string   `xml:",chardata"`
//}
//
//var Mechanisms = new(mechanismsElement)

type AuthElement struct {
	XMLName   xml.Name `xml:"auth"`
	Mechanism string   `xml:"mechanism,attr"`
	Data      string   `xml:",chardata"`
	elements.InnerElements
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

var ElementFactory = elements.NewFactory()

func init() {
	//features.List.AddElement(Mechanisms)
	features.ElementFactory.AddConstructor("urn:ietf:params:xml:ns:xmpp-sasl auth", func() elements.Element {
		return &AuthElement{InnerElements: elements.InnerElements{ElementFactory: ElementFactory}}
	})
}
