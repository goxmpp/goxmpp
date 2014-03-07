package features

import (
	"errors"

	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
)

type Container struct {
	*elements.InnerElements
}

func Loop(stream *stream.Stream) error {
	for stream.Opened && Tree.IsRequiredFor(stream) {
		stream.WriteElement(Tree.CopyIfAvailable(stream))
		e, _ := stream.ReadElement()
		if feature_handler, ok := e.(Handler); ok {
			if err := feature_handler.Handle(stream); err != nil {
				return err
			}
		} else {
			return errors.New("Non-handler element received.")
		}
	}
	return nil
}

var Tree = NewContainer()

type AccessControllable interface {
	CopyIfAvailable(*stream.Stream) elements.Element
}

type AccessController interface {
	IsRequiredFor(*stream.Stream) bool
}

type Handler interface {
	Handle(*stream.Stream) error
}

func NewContainer() *Container {
	return &Container{
		InnerElements: elements.NewInnerElements(nil),
	}
}

func (self *Container) CopyAvailableFeatures(stream *stream.Stream, dest *Container) {
	for _, feature := range self.Elements {
		if feature_element, ok := feature.(AccessControllable); ok {
			dest.AddElement(feature_element.CopyIfAvailable(stream))
		} else {
			dest.AddElement(feature)
		}
	}
}

func (self *Container) HasFeaturesRequiredFor(stream *stream.Stream) bool {
	for _, feature := range self.Elements {
		if feature_element, ok := feature.(AccessController); ok && feature_element.IsRequiredFor(stream) {
			return true
		}
	}
	return false
}

func (self *Container) IsRequiredFor(stream *stream.Stream) bool {
	return self.HasFeaturesRequiredFor(stream)
}

func (self *Container) CopyIfAvailable(stream *stream.Stream) elements.Element {
	e := NewContainer()
	self.CopyAvailableFeatures(stream, e)
	return e
}
