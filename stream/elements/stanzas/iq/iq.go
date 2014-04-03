package iq

import (
	"encoding/xml"
	"log"
)

import (
	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/stanzas"
)

const (
	STREAM_NODE = "iq"
)

func init() {
	stream.StreamFactory.AddConstructor(" "+STREAM_NODE, func() elements.Element {
		return NewIQElement()
	})
}

func NewIQElement() *IQElement {
	return &IQElement{InnerElements: elements.NewInnerElements(IQFactory)}
}

var IQFactory = elements.NewFactory()

type IQElement struct {
	XMLName xml.Name `xml:"iq"`
	stanzas.Base
	*elements.InnerElements
}

func (iq *IQElement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	iq.XMLName = start.Name

	iq.SetFromStartElement(start)

	return iq.UnmarshalInnerElements(d, start.End())
}

func (self *IQElement) Handle(stream *stream.Stream) error {
	log.Printf("Handling IQ: from = %#v, to = %#v\n", self.From, self.To)
	// TODO(dotdoom): the next line is for debugging only
	stream.Opened = false
	return nil
}
