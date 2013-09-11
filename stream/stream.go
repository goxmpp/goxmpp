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
	decodeStream(sw)
}

// This is an entry point for decode stanzas
func NextStanza(sw *Wrapper) {
	decodeStream(sw)
}

func decodeStream(sw *Wrapper) {
	for token, err := sw.StreamDecoder.Token(); err == nil; token, err = sw.StreamDecoder.Token() {
		switch element := token.(type) {
		case *xml.StartElement:
			handler, err := HandlerRegistrator.GetHandler(element.Name.Space + " " + element.Name.Local)
			if err != nil {
				// TODO: do logging here
				continue
			}

			if decode_err := sw.StreamDecoder.DecodeElement(handler, element); decode_err != nil {
				// TODO: do logging here
				continue
			}

			// All further processing goes in stanza hadler which may generate some output
			handler.HandleElement(sw)
		}
	}
}
