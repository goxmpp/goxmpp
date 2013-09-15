package message

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/stanza"

func init() {
	stream.GlobalElementFactory.AddConstructor(" message", func() stream.Element {
		return &Message{InnerXML: stream.InnerXML{ElementFactory: ElementFactory}}
	})

	ElementFactory.AddConstructor(" body", func() stream.Element {
		return &Body{}
	})

	ElementFactory.AddConstructor("*", func() stream.Element {
		return &InnerXML{}
	})
}

var ElementFactory = stream.NewElementFactory()

type InnerXML struct {
	XMLName xml.Name
	Text    string `xml:",innerxml"`
}

type Body struct {
	XMLName xml.Name `xml:"body"`
	Body    string   `xml:",chardata"`
}

type Message struct {
	XMLName xml.Name `xml:"message"`
	stanza.BaseStanza
	stream.InnerXML
}
