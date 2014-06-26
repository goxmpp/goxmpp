package message

import (
	"encoding/xml"

	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/stanzas"
	"github.com/goxmpp/xtream"
)

var messageXMLName = xml.Name{Local: "message"}
var bodyXMLName = xml.Name{Local: "body"}

func init() {
	xtream.NodeFactory.Add(func() xtream.Element {
		return NewMessageElement()
	}, stream.StreamXMLName, messageXMLName)

	xtream.NodeFactory.Add(func() xtream.Element {
		return &Body{}
	}, messageXMLName, bodyXMLName)
}

func NewMessageElement() *MessageElement {
	return &MessageElement{InnerElements: xtream.NewElements(&messageXMLName)}
}

type Body struct {
	XMLName xml.Name `xml:"body"`
	Body    string   `xml:",innerxml"`
}

type MessageElement struct {
	XMLName xml.Name `xml:"message"`
	stanzas.Base
	xtream.InnerElements `xml:",any"`
}
