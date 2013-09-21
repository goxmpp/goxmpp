package features

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/connection"
	"github.com/dotdoom/goxmpp/stream/elements"
)

var List = new(Features)
var Factory = elements.NewFactory()

type Entry interface {
	CopyIfAvailable(connection.State) interface{}
	IsRequiredFor(connection.State) bool
}

type Reactor interface {
	React(*connection.Connection)
}

type Elements struct {
	elements.Elements
}

func (self *Elements) CopyAvailableFeatures(fs connection.State, dest elements.ElementsAdder) elements.ElementsAdder {
	for _, feature := range self.Elements.Elements {
		if feature_entry, ok := feature.(Entry); ok {
			dest.AddElement(feature_entry.CopyIfAvailable(fs))
		} else {
			dest.AddElement(feature)
		}
	}
	return dest
}

func (self *Elements) HasFeaturesRequiredFor(fs connection.State) bool {
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

func (self *Features) IsRequiredFor(fs connection.State) bool {
	return self.HasFeaturesRequiredFor(fs)
}

func (self *Features) CopyIfAvailable(fs connection.State) interface{} {
	return self.CopyAvailableFeatures(fs, new(Features))
}
