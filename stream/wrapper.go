package stream

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/decoder"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
	"io"
)

type Wrapper struct {
	rw             io.ReadWriter
	StreamEncoder  *xml.Encoder
	StreamDecoder  *xml.Decoder
	ElementFactory elements.Factory
	FeatureSet     *features.FeaturesElement
	InnerDecoder   *decoder.InnerDecoder
	State          features.State
}

func (self *Wrapper) SetIO(rw io.ReadWriter) {
	self.rw = rw
	self.StreamEncoder = xml.NewEncoder(rw)
	self.StreamDecoder = xml.NewDecoder(rw)
}

func NewWrapper(rw io.ReadWriter) *Wrapper {
	w := &Wrapper{
		InnerDecoder:   decoder.NewInnerDecoder(),
		ElementFactory: features.Factory,
		FeatureSet:     features.List,
	}
	w.SetIO(rw)
	return w
}

func (self *Wrapper) NextElement() elements.Element {
	var element elements.Element

	elements.UnmarshalSiblingElements(self.StreamDecoder, self.ElementFactory, func(e elements.Element) bool {
		element = elements.ParseElement(e, self.InnerDecoder)
		return false
	})

	return element
}

func (self *Wrapper) FeaturesLoop() {
	for self.FeatureSet.IsRequiredFor(self.State) {
		self.StreamEncoder.Encode(self.FeatureSet.CopyIfAvailable(self.State))
		self.NextElement().(features.Reactor).React(self.State, self)
	}
}
