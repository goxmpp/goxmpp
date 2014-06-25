package iq

import (
	"encoding/xml"
	"errors"
	"log"
)

import (
	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/elements"
	"github.com/goxmpp/goxmpp/stream/elements/stanzas"
)

func init() {
	stream.StreamFactory.AddConstructor(func() elements.Element {
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

type Handler interface {
	Handle(*IQElement, *stream.Stream) error
}

func (self *IQElement) Handle(stream *stream.Stream) error {
	log.Printf("Handling IQ: from = %#v, to = %#v\n", self.From, self.To)

	match := false
	for _, element := range self.Elements() {
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
		// TODO(goxmpp): return more specific error class so it can be intercepted outside
		return errors.New("No inner elements handle this IQ")
	}
}
