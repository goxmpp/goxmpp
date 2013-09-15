package iq

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/elements"
import "github.com/dotdoom/goxmpp/stream/elements/stanza"

const (
	STREAD_NODE = "iq"
)

func init() {
	stream.GlobalElementFactory.AddConstructor(" "+STREAD_NODE, func() elements.Element {
		return &IQ{InnerXML: stream.InnerXML{ElementFactory: ElementFactory}}
	})
}

var ElementFactory = elements.NewFactory()

type IQ struct {
	XMLName xml.Name `xml:"iq"`
	stanza.BaseStanza
	stream.InnerXML
}
