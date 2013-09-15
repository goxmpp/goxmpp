package auth

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type Mechanisms struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	features.InnerElements
}

func (self *Mechanisms) IsRequiredFor(sw *stream.Wrapper) bool {
	return sw.State["authenticated"] == nil
}

func (self *Mechanisms) CopyIfAvailable(sw *stream.Wrapper) interface{} {
	if self.IsRequiredFor(sw) {
		return self.CopyAvailableInnerFeatures(sw, new(Mechanisms))
	} else {
		return nil
	}
}

type Mechanism struct {
	XMLName xml.Name `xml:"mechanism"`
	Name    string   `xml:",chardata"`
	features.InnerElements
}

func init() {
	features.GlobalFeaturesList.AddInnerElement(Mechanisms{})
}
