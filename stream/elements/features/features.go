package features

import (
	"encoding/xml"
	"errors"
	"log"

	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/elements"
)

type FeaturesElement struct {
	XMLName xml.Name `xml:"stream:features"`
	*Container
}

func NewFeaturesElement() *FeaturesElement {
	return &FeaturesElement{
		Container: NewContainer(),
	}
}

var Tree = NewFeaturesElement()

type Handler interface {
	Handle(*stream.Stream) error
}

func Loop(stream *stream.Stream) error {
	log.Println("entering features loop")
	for stream.Opened && Tree.IsRequiredFor(stream) {
		if err := stream.WriteElement(Tree.CopyIfAvailable(stream)); err != nil {
			return err
		}
		e, err := stream.ReadElement()
		if err != nil {
			return err
		}
		if feature_handler, ok := e.(Handler); ok {
			log.Println("calling feature handler")
			if err := feature_handler.Handle(stream); err != nil {
				return err
			}
			log.Println("feature handler completed")
		} else {
			return errors.New("Non-handler element received.")
		}

		if stream.ReOpen {
			stream.ReadSentOpen()
		}
	}
	log.Println("stream closed or no required features")
	return nil
}

func (self *FeaturesElement) CopyIfAvailable(stream *stream.Stream) elements.Element {
	e := NewFeaturesElement()
	self.CopyAvailableFeatures(stream, e.Container)
	return e
}
