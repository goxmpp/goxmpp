package features

import (
	"encoding/xml"

	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/xtream"
)

type Options interface{}
type FeatureHandler interface {
	Handle(*stream.Stream, Options) error
}

type BasicFeature interface {
	NewHandler() FeatureHandler
}

type Feature struct {
	name           string
	featureElement BasicFeature
	handlerElement FeatureHandler
	required       bool
}

func NewFeature(name string, felement BasicFeature, required bool) *Feature {
	return &Feature{name: name, featureElement: felement, required: required}
}

func (fw *Feature) Required() bool {
	return fw.required
}

func (fw *Feature) Name() string {
	return fw.name
}

func (fw *Feature) InitHandler() xtream.Element {
	fw.handlerElement = fw.featureElement.NewHandler()
	return fw
}

func (fw *Feature) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	return e.Encode(fw.featureElement)
}

func (fw *Feature) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return d.DecodeElement(fw.handlerElement, &start)
}

func (fw *Feature) Handle(strm *stream.Stream, opts Options) error {
	if err := fw.handlerElement.Handle(strm, opts); err != nil {
		return err
	}

	for _, dep := range strm.DependancyGraph().Get(fw.name) {
		if ffe := FeatureFactory.Get(dep); ffe != nil {
			f := ffe.Constructor(nil) // Get config for feature from stream
			strm.ElementFactory.AddNamed(
				func() xtream.Element { return f.InitHandler() },
				ffe.Parent,
				ffe.Name,
			)
			strm.AddFeature(f)
		}
	}

	strm.RemoveFeature(fw.name)

	return nil
}
