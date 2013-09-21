package auth

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type mechanisms struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	features.Elements
}

func (self *mechanisms) IsRequiredFor(fs features.FeatureState) bool {
	return fs["authenticated"] == nil
}

func (self *mechanisms) CopyIfAvailable(fs features.FeatureState) interface{} {
	if self.IsRequiredFor(fs) {
		return self.CopyAvailableFeatures(fs, new(mechanisms))
	}
	return nil
}

var Mechanisms = new(mechanisms)

type Mechanism struct {
	XMLName xml.Name `xml:"mechanism"`
	Name    string   `xml:",chardata"`
	features.Elements
}

func init() {
	features.GlobalFeaturesList.AddElement(Mechanisms)
}
