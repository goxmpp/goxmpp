package features

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements"
)

var GlobalFeaturesList = new(Features)

// TODO(artem): move to another package?
type FeatureState map[string]interface{}

type Entry interface {
	CopyIfAvailable(FeatureState) interface{}
	IsRequiredFor(FeatureState) bool
}

type Elements struct {
	elements.Elements
}

func (self *Elements) CopyAvailableFeatures(fs FeatureState, dest elements.ElementsAdder) elements.ElementsAdder {
	for _, feature := range self.Elements.Elements {
		self.AddElement(feature.(Entry).CopyIfAvailable(fs))
	}
	return self
}

func (self *Elements) HasFeaturesRequiredFor(fs FeatureState) bool {
	for _, feature := range self.Elements.Elements {
		if feature.(Entry).IsRequiredFor(fs) {
			return true
		}
	}
	return false
}

type Features struct {
	XMLName xml.Name `xml:"stream:features"`
	Elements
}
