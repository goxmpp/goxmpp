package stream

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/decoder"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
	"io"
)

type Wrapper struct {
	rwStream       io.ReadWriter
	StreamEncoder  *xml.Encoder
	StreamDecoder  *xml.Decoder
	InnerDecoder   *decoder.InnerDecoder
	ElementFactory elements.Factory
	FeatureSet     *features.Features
	State          features.FeatureState
}

func NewWrapper(rw io.ReadWriter) *Wrapper {
	return &Wrapper{
		rwStream:       rw,
		StreamEncoder:  xml.NewEncoder(rw),
		StreamDecoder:  xml.NewDecoder(rw),
		InnerDecoder:   decoder.NewInnerDecoder(),
		State:          features.FeatureState{},
		ElementFactory: features.Factory,
		FeatureSet:     features.List,
	}
}

func (self *Wrapper) SwapIOStream(rw io.ReadWriter) {
	self.rwStream = rw
	self.StreamEncoder = xml.NewEncoder(rw)
	self.StreamDecoder = xml.NewDecoder(rw)
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
		self.NextElement().(features.Reactor).React()
		break
	}
}
