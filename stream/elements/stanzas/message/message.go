package message

import "encoding/xml"

import (
	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/stanzas"
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

func (msg *MessageElement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	msg.XMLName = start.Name

	msg.SetFromStartElement(start)

	return msg.UnmarshalInnerElements(d, start.End())
}
