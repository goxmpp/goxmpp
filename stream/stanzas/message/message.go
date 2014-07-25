package message

import (
	"encoding/xml"

	"github.com/goxmpp/goxmpp/stream/stanzas"
	"github.com/goxmpp/xtream"
)

var messageXMLName = xml.Name{Local: "message"}
var bodyXMLName = xml.Name{Local: "body"}

func init() {
	xtream.NodeFactory.Add(func() xtream.Element {
		return NewMessageElement()
	})

	xtream.NodeFactory.Add(func() xtream.Element {
		return &Body{}
	})
}

func NewMessageElement() *MessageElement {
	return &MessageElement{InnerElements: xtream.NewElements()}
}

type Body struct {
	XMLName xml.Name `xml:"body" parent:"message"`
	Body    string   `xml:",innerxml"`
}

type MessageElement struct {
	XMLName xml.Name `xml:"message" parent:"stream:stream"`
	stanzas.Base
	xtream.InnerElements `xml:",any"`
}
