package features

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements"
	"io"
)

var List = new(Features)
var Factory = elements.NewFactory()

type State map[string]interface{}

type Entry interface {
	CopyIfAvailable(State) interface{}
	IsRequiredFor(State) bool
}

type SuperInterface interface {
	SetIO(io.ReadWriter)
	NextElement() elements.Element
}

type Reactor interface {
	React(State, SuperInterface)
}

type Elements struct {
	elements.Elements
}

func (self *Elements) CopyAvailableFeatures(fs State, dest elements.ElementsAdder) elements.ElementsAdder {
	for _, feature := range self.Elements.Elements {
		if feature_entry, ok := feature.(Entry); ok {
			dest.AddElement(feature_entry.CopyIfAvailable(fs))
		} else {
			dest.AddElement(feature)
		}
	}
	return dest
}

func (self *Elements) HasFeaturesRequiredFor(fs State) bool {
	for _, feature := range self.Elements.Elements {
		if feature.(Entry).IsRequiredFor(fs) {
			return true
		}
	}
	return false
}

type Features struct {
	XMLName xml.Name `xml:"stream:features"`
	Elements
}

func (self *Features) IsRequiredFor(fs State) bool {
	return self.HasFeaturesRequiredFor(fs)
}

func (self *Features) CopyIfAvailable(fs State) interface{} {
	return self.CopyAvailableFeatures(fs, new(Features))
}
