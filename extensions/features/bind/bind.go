package bind

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type bind struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	features.Elements
}

func (self *bind) IsRequiredFor(fs features.FeatureState) bool {
	return fs["bound"] == nil
}

func (self *bind) CopyIfAvailable(fs features.FeatureState) interface{} {
	if self.IsRequiredFor(fs) {
		return self.CopyAvailableFeatures(fs, new(bind))
	}
	return nil
}

func init() {
	features.GlobalFeaturesList.AddElement(new(bind))
}
