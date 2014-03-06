package stream

import (
	"encoding/xml"
	"errors"
	"io"

	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type Stream struct {
	XMLName          xml.Name
	ID               string `xml:"id,attr"`
	From             string `xml:"from,attr"` // This holds user JID after auth.
	To               string `xml:"to,attr"`
	Version          string `xml:"version,attr"`
	DefaultNamespace string `xml:"-"`
	Opened           bool   `xml:"-"`
	Features         features.State
	Connection
	elements.ElementFactory
}

var Factory = elements.NewElementFactory()

func NewStream(rw io.ReadWriter) *Stream {
	st := &Stream{
		ElementFactory: Factory,
	}
	st.SetRW(rw)
	return st
}

func (self *Stream) Parse(_ *xml.Decoder, start *xml.StartElement) error {
	self.XMLName = start.Name
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "to":
			self.To = attr.Value
		case "from":
			self.From = attr.Value
		case "version":
			self.Version = attr.Value
		}
	}
	return nil
}

// TODO(artem): refactor
func (self *Stream) Compose(_ *xml.Encoder, w io.Writer) error {
	data := xml.Header

	data += "<stream:" + self.XMLName.Local + " xmlns='" + self.DefaultNamespace + "' xmlns:stream='" + self.XMLName.Space + "'"
	if self.ID != "" {
		data += " id='" + self.ID + "'"
	}
	if self.From != "" {
		data += " from='" + self.From + "'"
	}
	if self.To != "" {
		data += " to='" + self.To + "'"
	}
	if self.Version != "" {
		data += " version='" + self.Version + "'"
	}
	data += ">"

	_, err := io.WriteString(w, data)
	return err
}

func (self *Stream) WriteElement(element elements.Element) error {
	return self.streamEncoder.Encode(element)
}

func (self *Stream) ReadElement() (elements.Element, error) {
	var element elements.Element
	var err error
	for token, err := self.streamDecoder.Token(); err == nil; token, err = self.streamDecoder.Token() {
		if start, ok := token.(xml.StartElement); ok {
			element, err = Factory.DecodeElement(self.streamDecoder, &start)
			if err != nil {
				return nil, err
			}
		}
	}

	return element, err
}

func (self *Stream) FeaturesLoop() error {
	for self.Opened && features.Features.IsRequiredFor(&self.Features) {
		self.WriteElement(features.Features.CopyIfAvailable(&self.Features))
		e, _ := self.ReadElement()
		if feature_handler, ok := e.(features.Handler); ok {
			feature_handler.Handle(&self.Features)
		} else {
			return errors.New("Non-handler element received.")
		}
	}
	return nil
}
