package features

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/stream/elements"
)

var ElementFactory = elements.NewFactory()

// stream:features element
type FeaturesElement struct {
	XMLName xml.Name `xml:"stream:features"`
	Elements
	elements.InnerElements
}

func NewFeaturesElement() *FeaturesElement {
	return &FeaturesElement{
		InnerElements: elements.InnerElements{ElementFactory: ElementFactory},
	}
}

// Helper methods for feature element containing sub-elements as a slice.
type Elements struct {
	Elements []interface{}
}

func (self *Elements) AddElement(element interface{}) {
	if element != nil {
		self.Elements = append(self.Elements, element)
	}
}

func (self *Elements) CopyAvailableFeatures(fs interface{}, dest *Elements) {
	for _, feature := range self.Elements {
		if feature_element, ok := feature.(interface {
			CopyIfAvailable(interface{}) interface{}
		}); ok {
			dest.AddElement(feature_element.CopyIfAvailable(fs))
		} else {
			dest.AddElement(feature)
		}
	}
}

func (self *Elements) HasFeaturesRequiredFor(fs interface{}) bool {
	for _, feature := range self.Elements {
		if feature_element, ok := feature.(interface {
			IsRequiredFor(interface{}) bool
		}); ok && feature_element.IsRequiredFor(fs) {
			return true
		}
	}
	return false
}

func (self *FeaturesElement) IsRequiredFor(fs interface{}) bool {
	return self.HasFeaturesRequiredFor(fs)
}

func (self *FeaturesElement) CopyIfAvailable(fs interface{}) interface{} {
	e := new(FeaturesElement)
	self.CopyAvailableFeatures(fs, &e.Elements)
	return e
}
