package features

import "github.com/dotdoom/goxmpp/stream/elements"

type Container struct {
	*elements.InnerElements
}

type AccessControllable interface {
	CopyIfAvailable(*State) elements.Element
	IsRequiredFor(*State) bool
}

func NewContainer(factory elements.ElementFactory) *Container {
	return &Container{
		InnerElements: elements.NewInnerElements(factory),
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

func (self *FeaturesElement) IsRequiredFor(fs *State) bool {
	return self.HasFeaturesRequiredFor(fs)
}

func (self *FeaturesElement) CopyIfAvailable(fs *State) interface{} {
	e := NewContainer(nil)
	self.CopyAvailableFeatures(fs, e)
	return e
}
