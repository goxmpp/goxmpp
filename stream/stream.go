package stream

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements"
)

type Stream struct {
	XMLName xml.Name
	ID      string
	From    string
	To      string
	Version string
}

// An entry point for decoding elements in response to features
// announcment from server, before session is opened
func handleFeature(sw *Wrapper) {
	elements.UnmarshalSiblingElements(sw.StreamDecoder, sw.ElementFactory, func(element elements.Element) bool {
		// TODO: need to check sw for the state when all required features processed and exit the loop
		elements.UnmarshalElement(element, sw.InnerDecoder)
		return true
	})
}

// This is an entry point for decode stanzas
func NextStanza(sw *Wrapper) elements.Element {
	var stanza elements.Element

	elements.UnmarshalSiblingElements(sw.StreamDecoder, sw.ElementFactory, func(element elements.Element) bool {
		stanza = elements.UnmarshalElement(element, sw.InnerDecoder)
		return false // We need process stanzas one by one
	})

	return stanza
}
