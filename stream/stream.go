package stream

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements"
)

var GlobalStanzasFactory = elements.NewFactory()
var GlobalFeaturesFactory = elements.NewFactory()

type Stream struct {
	XMLName xml.Name
	ID      string
	From    string
	To      string
	Version string
}

type Element interface{}

// An entry point for decoding elements in response to features
// announcment from server, before session is opened
func handleFeature(sw *Wrapper) {
	unmarshalSiblingElements(sw.StreamDecoder, sw.ElementFactory, func(element Element) bool {
		// TODO: need to check sw for the state when all required features processed and exit the loop
		unmarshalElement(element, sw.InnerDecoder)
		return true
	})
}

// This is an entry point for decode stanzas
func NextStanza(sw *Wrapper) Element {
	var stanza Element

	unmarshalSiblingElements(sw.StreamDecoder, sw.ElementFactory, func(element Element) bool {
		stanza = unmarshalElement(element, sw.InnerDecoder)
		return false // We need process stanzas one by one
	})

	return stanza
}
