package features

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream/elements/stanza"

type FeatureEntry interface {
	CopyIfAvailable(*StreamWrapper) interface{}
	IsRequiredFor(*StreamWrapper) bool
}

type FeaturesInnerElement struct {
	stanza.InnerElements
}

func (self *FeatureEntry) CopyInnerElements(sw *StreamWrapper, dest stream.InnerElementAdder) stream.InnerElementAdder {
	for _, feature := range self.InnerElements {
		self.AddInnerElement(feature.CopyIfAvailable(sw))
	}
	return sf
}

type Features struct {
	XMLName xml.Name `xml:"stream:features"`
	stream.InnerElements
}
