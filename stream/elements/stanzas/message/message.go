package message

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/elements"
import "github.com/dotdoom/goxmpp/stream/elements/stanzas"

func init() {
	stream.Factory.AddConstructor("message", func() elements.Element {
		return NewMessage()
	})

	ElementFactory.AddConstructor("body", func() elements.Element {
		return &Body{}
	})
}

func NewMessage() *Message {
	return &Message{InnerElements: elements.NewInnerElements(ElementFactory)}
}

var ElementFactory = elements.NewFactory()

type Body struct {
	XMLName xml.Name `xml:"body"`
	Body    string   `xml:",innerxml"`
}

type Message struct {
	XMLName xml.Name `xml:"message"`
	stanzas.Base
	*elements.InnerElements
}

func (msg *Message) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	msg.XMLName = start.Name

	msg.SetFromStartElement(start)

	return msg.HandlerInnerElements(d, start.End())
}
