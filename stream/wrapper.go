package stream

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/decoder"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
	"io"
)

type Connection struct {
	rw             io.ReadWriter
	streamEncoder  *xml.Encoder
	streamDecoder  *xml.Decoder
	ElementFactory elements.Factory
	InnerDecoder   decoder.InnerDecoder
}

func (self *Connection) SetIO(rw io.ReadWriter) io.ReadWriter {
	previous_rw := self.rw
	self.rw = rw
	self.streamEncoder = xml.NewEncoder(rw)
	self.streamDecoder = xml.NewDecoder(rw)
	return previous_rw
}

func (self *Connection) ReadElement() (element elements.Parsable, fatal error) {
	elements.ParseSiblingElements(self.streamDecoder, self.ElementFactory, func(e elements.Parsable) bool {
		element = e
		if e, ok := e.(elements.ParsableElements); ok {
			fatal = e.ParseElements(&self.InnerDecoder)
		}
		return false
	})

	return element, fatal
}

func (self *Connection) WriteElement(element elements.Composable) error {
	return element.Compose(self.streamEncoder, self.rw)
}

func (self *Connection) WritePrompt(string) error {
	return nil
}

func (self *Connection) FeaturesLoop(fe *features.FeaturesElement, state interface{}) {
	for state.(interface {
		Opened() bool
	}).Opened() && fe.IsRequiredFor(state) {
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
