package message

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"

type Message struct {
	XMLName xml.Name `xml:"message"`
	stream.Stanza
	Body string `xml:"body,omitempty"`
}
