package stream

import "github.com/dotdoom/goxmpp/extensions"

var HandlerRegistrator = extensions.NewHandlerRegistrator()

type Stream struct {
	XMLName xml.Name
	ID      string
	From    string
	To      string
	Version string
}

// An entry point for decoding elements in response to features
// anouncment from server, before sission is opened
func HandleFeature(sw *Wrapper) {
	decode_stream(sw)
}

// This is an entry point for decode stanzas
func NextStanza(sw *Wrapper) {
	decode_stream(sw)
}

func decode_stream(sw *Wrapper) {
	for token, err := sw.StreamDecoder.Token(); err == nil; token, err = sw.StreamDecoder.Token() {
		switch element, real_type := token.(type); real_type {
		case xml.StartElement:
			handler := HandlerRegistrator.GetHandler(elemnt.Name.Space + " " + element.Name.Local)
			if decode_err := sw.StreamDecoder.DecodeElement(handler, element); decode_err != nil {
				// TODO: do logging here
			}

			// All further processing goes in stanza hadler which may generate some output
			handler.HandleElement(sw)
		}
	}
}
