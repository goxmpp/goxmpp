package iq

import (
	"encoding/xml"
	"errors"
	"log"

	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/stanzas"
	"github.com/goxmpp/xtream"
)

var IQXMLName = xml.Name{Local: "iq"}

func init() {
	xtream.NodeFactory.Add(func() xtream.Element {
		return NewIQElement()
	}, stream.StreamXMLName, IQXMLName)
}

func NewIQElement() *IQElement {
	return &IQElement{InnerElements: xtream.NewElements(&IQXMLName)}
}

type IQElement struct {
	XMLName xml.Name `xml:"iq"`
	stanzas.Base
	xtream.InnerElements `xml:",any"`
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
