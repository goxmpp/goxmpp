package iq

import (
	"encoding/xml"
	"errors"
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

type Handler interface {
	Handle(*IQElement, *stream.Stream) error
}

func (self *IQElement) Handle(stream *stream.Stream) error {
	log.Printf("Handling IQ: from = %#v, to = %#v\n", self.From, self.To)

	match := false
	for _, element := range self.Elements {
		if handler, ok := element.(Handler); ok {
			match = true
			if err := handler.Handle(self, stream); err != nil {
				return err
			}
		}
	}

	if match {
		return nil
	} else {
		// TODO(dotdoom): return more specific error class so it can be intercepted outside
		return errors.New("No inner elements handle this IQ")
	}
}
