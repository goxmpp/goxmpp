package features

import (
	"encoding/xml"
	"log"

	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/xtream"
)

type Options interface{}
type FeatureHandler interface {
	Handle(stream.ServerStream, Options) error
}

type BasicFeature interface {
	NewHandler() FeatureHandler
}

type Feature struct {
	name           string
	featureElement BasicFeature
	handlerElement FeatureHandler
	required       bool
	config         Options
}

func NewFeature(name string, felement BasicFeature, required bool, conf Options) *Feature {
	return &Feature{name: name, featureElement: felement, required: required, config: conf}
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

func (fw *Feature) Handle(strm stream.ServerStream, opts Options) error {
	if opts == nil {
		opts = fw.config
	}

	if err := fw.handlerElement.Handle(strm, opts); err != nil {
		return err
	}

	EnableStreamFeatures(strm, fw.name)

	strm.RemoveFeature(fw.name)

	return nil
}

func EnableStreamFeatures(s stream.ServerStream, name string) {
	for _, fname := range s.DependencyGraph().Get(name) {
		fe := FeatureFactory.Get(fname)

		conf, err := fe.GetConfig(s.Config()[fname])
		if err != nil {
			log.Printf("goxmpp: unable to handle config for feature %s: %s", fname, err)
			continue
		}

		feature := fe.Constructor(conf)

		s.AddNamed(
			func() xtream.Element { return feature.InitHandler() },
			fe.Parent,
			fe.Name,
		)
		s.AddFeature(feature)
	}
}
