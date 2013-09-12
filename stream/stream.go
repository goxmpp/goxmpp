package stream

import (
	"encoding/xml"
)

var HandlerRegistrator = NewElementGeneratorRegistrator()

type Stream struct {
	XMLName xml.Name
	ID      string
	From    string
	To      string
	Version string
}

type Element interface{}

// An entry point for decoding elements in response to features
// anouncment from server, before sission is opened
func handleFeature(sw *Wrapper) {
	processStreamElements(sw.StreamDecoder, HandlerRegistrator, func(handler Element) bool {
		// TODO: need to check sw for the state when all required features processed and exit the loop
		unmarshalStreamElement(handler, sw)
		return true
	})
}

// This is an entry point for decode stanzas
func NextStanza(sw *Wrapper) Element {
	var stanza Element

	processStreamElements(sw.StreamDecoder, HandlerRegistrator, func(handler Element) bool {
		stanza = unmarshalStreamElement(handler, sw)
		return false // We need process stanzas one by one
	})

	return stanza
}
