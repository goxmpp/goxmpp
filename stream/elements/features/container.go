package features

import (
	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/elements"
)

type Container struct {
	*elements.InnerElements
}

type AccessControllable interface {
	CopyIfAvailable(*stream.Stream) elements.Element
}

type AccessController interface {
	IsRequiredFor(*stream.Stream) bool
}

func NewContainer() *Container {
	return &Container{
		InnerElements: elements.NewInnerElements(nil),
	}
}

func (self *Container) CopyAvailableFeatures(stream *stream.Stream, dest *Container) {
	for _, feature := range self.Elements() {
		if feature_element, ok := feature.(AccessControllable); ok {
			dest.AddElement(feature_element.CopyIfAvailable(stream))
		} else {
			dest.AddElement(feature)
		}
	}
}

func (self *Container) HasFeaturesRequiredFor(stream *stream.Stream) bool {
	for _, feature := range self.Elements() {
		if feature_element, ok := feature.(AccessController); ok && feature_element.IsRequiredFor(stream) {
			return true
		}
	}
	return false
}

func (self *Container) IsRequiredFor(stream *stream.Stream) bool {
	return self.HasFeaturesRequiredFor(stream)
}
