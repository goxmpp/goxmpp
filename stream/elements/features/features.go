package features

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/stream/elements"
)

var Factory = elements.NewElementFactory()

type FeaturesElement struct {
	XMLName xml.Name `xml:"stream:features"`
	*Container
}

func NewFeaturesElement() *FeaturesElement {
	return &FeaturesElement{
		Container: NewContainer(nil),
	}
}
