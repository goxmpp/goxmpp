package features

import (
	"encoding/xml"

	"github.com/goxmpp/xtream"
)

type FeatureHandler interface {
	Handle(FeatureContainable) error
}

type Feature struct {
	Name           string
	featureElement xtream.Element
	handlerElement FeatureHandler
}

func NewFeature(name string, felement xtream.Element) *Feature {
	return &Feature{Name: name, featureElement: felement}
}

func (fw *Feature) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	return e.Encode(fw.featureElement)
}

func (fw *Feature) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return d.DecodeElement(fw.handlerElement, &start)
}

func (fw *Feature) Handle(strm FeatureContainable) error {
	if err := fw.handlerElement.Handle(strm); err != nil {
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
