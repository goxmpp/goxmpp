package iq

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/stanza"

const (
	STREAD_NODE = "iq"
)

func init() {
	stream.GlobalElementFactory.AddConstructor(" "+STREAD_NODE, func() stream.Element {
		return &IQ{InnerXML: stream.InnerXML{ElementFactory: ElementFactory}}
	})
}

var ElementFactory = stream.NewElementFactory()

type IQ struct {
	XMLName xml.Name `xml:"iq"`
	stanza.BaseStanza
	stream.InnerXML
}
