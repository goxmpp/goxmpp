package features

import "encoding/xml"

type FeatureContainable interface {
	AddFeature(*Feature)
	RemoveFeature(string)
	HasRequired() bool
}

type FeatureContainer struct {
	features     map[string]*Feature
	num_required int
}

func NewFeatureContainer() *FeatureContainer {
	return &FeatureContainer{features: make(map[string]*Feature)}
}

func (fc *FeatureContainer) HasRequired() bool {
	return fc.num_required > 0
}

func (fc *FeatureContainer) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start = xml.StartElement{Name: xml.Name{Local: "stream:features"}}
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	fs := make([]*Feature, 0, len(fc.features))
	for _, v := range fc.features {
		fs = append(fs, v)
	}

	if err := e.Encode(fs); err != nil {
		return err
	}

	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}

	return e.Flush()
}

func (fc *FeatureContainer) AddFeature(f *Feature) {
	if _, ok := fc.features[f.Name]; !ok {
		fc.features[f.Name] = f
		if f.Required {
			fc.num_required += 1
		}
	}
}

func (fc *FeatureContainer) RemoveFeature(name string) {
	if feature, ok := fc.features[name]; ok && feature.Required {
		fc.num_required -= 1
	}
	delete(fc.features, name)
}
