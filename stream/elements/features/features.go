package features

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/stream/elements"
)

var Factory = elements.NewElementFactory()

// stream:features element
type FeaturesElement struct {
	XMLName xml.Name `xml:"stream:features"`
	*Container
}

func NewFeaturesElement() *FeaturesElement {
	return &FeaturesElement{
		Container: NewContainer(Factory),
	}
}
