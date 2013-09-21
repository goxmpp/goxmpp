package iq

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream/elements"
import "github.com/dotdoom/goxmpp/stream/elements/stanzas"

const (
	STREAM_NODE = "iq"
)

func init() {
	elements.GlobalStanzasFactory.AddConstructor(" "+STREAM_NODE, func() elements.Element {
		return &IQ{UnmarshallableElements: elements.UnmarshallableElements{ElementFactory: ElementFactory}}
	})
}

var ElementFactory = elements.NewFactory()

type IQ struct {
	XMLName xml.Name `xml:"iq"`
	stanzas.Base
	elements.UnmarshallableElements
}
