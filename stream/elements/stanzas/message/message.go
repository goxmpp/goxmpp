package message

import "encoding/xml"

import (
	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/elements"
	"github.com/goxmpp/goxmpp/stream/elements/stanzas"
)

func init() {
	stream.StreamFactory.AddConstructor(func() elements.Element {
		return NewMessageElement()
	})

	MessageFactory.AddConstructor(func() elements.Element {
		return &Body{}
	})
}

func NewMessageElement() *MessageElement {
	return &MessageElement{InnerElements: elements.NewInnerElements(MessageFactory)}
}

var MessageFactory = elements.NewFactory()

type Body struct {
	XMLName xml.Name `xml:"body"`
	Body    string   `xml:",innerxml"`
}

type MessageElement struct {
	XMLName xml.Name `xml:"message"`
	stanzas.Base
	*elements.InnerElements
}
