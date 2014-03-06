package features

import "github.com/dotdoom/goxmpp/stream/elements"

type Container struct {
	*elements.InnerElements
}

var Features = NewContainer()

type AccessControllable interface {
	CopyIfAvailable(*State) elements.Element
	IsRequiredFor(*State) bool
}

type Handler interface {
	Handle(*State) bool
}

func NewContainer() *Container {
	return &Container{
		InnerElements: elements.NewInnerElements(nil),
	}
}

func (self *Container) CopyAvailableFeatures(fs *State, dest *Container) {
	for _, feature := range self.Elements {
		if feature_element, ok := feature.(AccessControllable); ok {
			dest.AddElement(feature_element.CopyIfAvailable(fs))
		} else {
			dest.AddElement(feature)
		}
	}
}

func (self *Container) HasFeaturesRequiredFor(fs *State) bool {
	for _, feature := range self.Elements {
		if feature_element, ok := feature.(AccessControllable); ok && feature_element.IsRequiredFor(fs) {
			return true
		}
	}
	return false
}

func (self *Container) IsRequiredFor(fs *State) bool {
	return self.HasFeaturesRequiredFor(fs)
}

func (self *Container) CopyIfAvailable(fs *State) elements.Element {
	e := NewContainer()
	self.CopyAvailableFeatures(fs, e)
	return e
}
