package presence

import (
	"encoding/xml"

	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/elements"
	"github.com/goxmpp/goxmpp/stream/elements/stanzas"
)

func init() {
	stream.StreamFactory.AddConstructor(func() elements.Element {
		return NewPresenceElement()
	})
}

func NewPresenceElement() *PresenceElement {
	return &PresenceElement{InnerElements: elements.NewInnerElements(PresenceFactory)}
}

var PresenceFactory = elements.NewFactory()

type PresenceElement struct {
	XMLName xml.Name `xml:"presence"`
	Show    string   `xml:"show,omitempty"`
	Status  string   `xml:"status"`
	stanzas.Base
	*elements.InnerElements
}
