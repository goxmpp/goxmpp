package stream

import (
	"github.com/dotdoom/goxmpp/stream/connection"
	"github.com/dotdoom/goxmpp/stream/decoder"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
	"io"
)

type Wrapper struct {
	ElementFactory elements.Factory
	FeatureSet     *features.Features
	InnerDecoder   *decoder.InnerDecoder
	*connection.Connection
}

func NewWrapper(rw io.ReadWriter) *Wrapper {
	return &Wrapper{
		Connection:     connection.NewConnection(rw),
		InnerDecoder:   decoder.NewInnerDecoder(),
		ElementFactory: features.Factory,
		FeatureSet:     features.List,
	}
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
		self.NextElement().(features.Reactor).React(self.Connection)
		break
	}
}
