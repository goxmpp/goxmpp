package message

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/stanza"

func init() {
	stream.HandlerRegistrator.Register(" message", func() stream.Element {
		return &Message{InnerXML: stream.InnerXML{Registrator: HandlerRegistrator}}
	})

	HandlerRegistrator.Register(" body", func() stream.Element {
		return &Body{}
	})

	HandlerRegistrator.Register("*", func() stream.Element {
		return &InnerXML{}
	})
}

var HandlerRegistrator = stream.NewElementGeneratorRegistrator()

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
