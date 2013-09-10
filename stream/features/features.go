package features

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"

type FeatureCopier interface {
	CopyIfAvailable(*StreamWrapper) interface{}
	IsRequiredFor(*StreamWrapper) bool
}

type FeatureEntry struct {
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
