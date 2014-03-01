package mechanism

import (
//	"bytes"
//	"encoding/base64"
//	"encoding/xml"
//	"github.com/dotdoom/goxmpp/extensions/features/auth"
//	"github.com/dotdoom/goxmpp/stream/elements"
//	"github.com/dotdoom/goxmpp/stream/elements/features"
//	"log"
)

//
//type PlainElement struct {
//	auth.MechanismElement
//}
//
//type SuccessElement struct {
//	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
//}
//
//func NewPlainElement() *PlainElement {
//	return &PlainElement{MechanismElement: auth.MechanismElement{Name: "PLAIN"}}
//}
//
//type ElementWriter interface {
//	WriteElement(elements.Element)
//}
//
//var usernamePasswordSeparator = []byte{0}
//
//func (self *PlainElement) Handle(a *auth.AuthElement, s features.State, sw interface{}) bool {
//	b, _ := base64.StdEncoding.DecodeString(a.Data)
//	user_password := bytes.Split(b, usernamePasswordSeparator)
//	log.Println("auth info:", string(user_password[1]), string(user_password[2]))
//	sw.(ElementWriter).WriteElement(&SuccessElement{})
//	s["authenticated"] = true
//	s["stream_opened"] = false
//	return true
//}
//
//func init() {
//	auth.Mechanisms.AddElement(NewPlainElement())
//}
