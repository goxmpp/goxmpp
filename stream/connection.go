package stream

import (
	"encoding/xml"
	"io"

	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type Connection struct {
	rw            io.ReadWriter
	streamEncoder *xml.Encoder
	streamDecoder *xml.Decoder
	stream        *Stream
}

type State struct {
	Opened bool
}

func NewConnection(rw io.ReadWriter) *Connection {
	self := &Connection{
		rw:            rw,
		streamEncoder: xml.NewEncoder(rw),
		streamDecoder: xml.NewDecoder(rw),
		stream:        NewStream(),
	}
	return self
}

func (self *Connection) ReadElement() (interface{}, error) {
	var element interface{}
	var err error
	for token, err := self.streamDecoder.Token(); err == nil; token, err = self.streamDecoder.Token() {
		if start, ok := token.(xml.StartElement); ok {
			element, err = self.stream.DecodeElemenet(self.streamDecoder, &start)
			if err != nil {
				return nil, err
			}
		}
	}

	return element, err
}

func (self *Connection) WriteElement(element elements.Element) error {
	return self.streamEncoder.Encode(element)
}

func (self *Connection) WritePrompt(string) error {
	return nil
}

func (self *Connection) FeaturesLoop(fe *features.FeaturesElement, state *features.State) {
	var connection_state *State
	state.Get(&connection_state)
	for connection_state.Opened && fe.IsRequiredFor(state) {
		self.WriteElement(fe.CopyIfAvailable(state))
		e, _ := self.ReadElement()
		if feature_handler, ok := e.(interface {
			HandleFeature(*Connection, interface{})
		}); ok {
			feature_handler.HandleFeature(self, state)
		} else {
			self.WritePrompt("Expecting required feature usage")
		}
	}
}
