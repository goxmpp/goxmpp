package stream

import (
	"encoding/xml"
)

var HandlerRegistrator = NewElementHandlerRegistrator()

type Stream struct {
	XMLName xml.Name
	ID      string
	From    string
	To      string
	Version string
}

// An entry point for decoding elements in response to features
// anouncment from server, before sission is opened
func handleFeature(sw *Wrapper) {
	processStreamElements(sw.StreamDecoder, HandlerRegistrator, func(handler ElementHandler) {
		handler.HandleElement(sw)
	})
}

// This is an entry point for decode stanzas
func NextStanza(sw *Wrapper) {
	processStreamElements(sw.StreamDecoder, HandlerRegistrator, func(handler ElementHandler) {
		handler.HandleElement(sw)
	})
}
