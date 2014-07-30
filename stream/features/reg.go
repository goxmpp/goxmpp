package features

import (
	"encoding/json"
	"encoding/xml"
)

type FeatureConstructor func(Options) *Feature
type FeatureFactoryElement struct {
	Constructor FeatureConstructor
	Config      interface{}
	Name        xml.Name
	Parent      xml.Name
	Wants       []string
}

func (ffe *FeatureFactoryElement) GetConfig(conf json.RawMessage) (interface{}, error) {
	if fn, ok := ffe.Config.(func() interface{}); ok {
		config := fn()
		if err := json.Unmarshal(conf, config); err != nil {
			return nil, err
		}
		ffe.Config = config
	}

	return ffe.Config, nil
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

	if len(ffe.Wants) == 0 {
		ffe.Wants = []string{"stream"}
	}
	DependencyGraph.Add(name, ffe.Wants...)
	ff.feature_cons[name] = ffe
}

func (ff *featureFactory) Get(name string) *FeatureFactoryElement {
	return ff.feature_cons[name]
}

func (ff *featureFactory) List() map[string]*FeatureFactoryElement {
	return ff.feature_cons
}
