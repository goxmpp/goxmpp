package features

import (
	"encoding/xml"

	"github.com/goxmpp/xtream"
)

type Options interface{}
type FeatureHandler interface {
	Handle(FeatureContainable, Options) error
}

type BasicFeature interface {
	NewHandler() FeatureHandler
}

type Feature struct {
	Name           string
	featureElement BasicFeature
	handlerElement FeatureHandler
}

func NewFeature(name string, felement BasicFeature) *Feature {
	return &Feature{Name: name, featureElement: felement}
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

func (fw *Feature) Handle(strm FeatureContainable, opts Options) error {
	if err := fw.handlerElement.Handle(strm, opts); err != nil {
		return err
	}

	/*for _, dep := range depgraph.ListNodes(fw.Name) {
		if feature := Features.Get(dep); feature != nil {
			strm.AddFeature(feature)
		}
	}*/

	strm.RemoveFeature(fw.Name)

	return nil
}
