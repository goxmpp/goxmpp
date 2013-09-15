package features

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements"
)

var GlobalFeaturesList = new(Features)

type FeatureState map[string]interface{}

type Entry interface {
	CopyIfAvailable(FeatureState) interface{}
	IsRequiredFor(FeatureState) bool
}

type InnerElements struct {
	elements.InnerElements
}

func (self *InnerElements) CopyAvailableInnerFeatures(fs FeatureState, dest elements.InnerElementsAdder) elements.InnerElementsAdder {
	for _, feature := range self.InnerElements.InnerElements {
		self.AddInnerElement(feature.(Entry).CopyIfAvailable(fs))
	}
	return self
}

func (self *InnerElements) HasInnerFeaturesRequiredFor(fs FeatureState) bool {
	for _, feature := range self.InnerElements.InnerElements {
		if feature.(Entry).IsRequiredFor(fs) {
			return true
		}
	}
	return false
}

type Features struct {
	XMLName xml.Name `xml:"stream:features"`
	InnerElements
}
