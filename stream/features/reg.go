package features

import "encoding/xml"

type FeatureConstructor func(Options) *Feature
type FeatureFactoryElement struct {
	Constructor FeatureConstructor
	Name        xml.Name
	Parent      xml.Name
}
type FF interface {
	Add(string, *FeatureFactoryElement)
	Get(string) *FeatureFactoryElement
	List() map[string]*FeatureFactoryElement
}

type featureFactory struct {
	feature_cons map[string]*FeatureFactoryElement
}

var FeatureFactory FF = newFactory()

func newFactory() *featureFactory {
	return &featureFactory{make(map[string]*FeatureFactoryElement)}
}

func (ff *featureFactory) Add(name string, ffe *FeatureFactoryElement) {
	if _, ok := ff.feature_cons[name]; ok {
		panic("Feature element already registered")
	}

	ff.feature_cons[name] = ffe
}

func (ff *featureFactory) Get(name string) *FeatureFactoryElement {
	return ff.feature_cons[name]
}

func (ff *featureFactory) List() map[string]*FeatureFactoryElement {
	return ff.feature_cons
}
