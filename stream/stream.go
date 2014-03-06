package stream

import (
	"encoding/xml"
	"io"

	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type Stream struct {
	XMLName          xml.Name
	ID               string `xml:"id,attr"`
	From             string `xml:"from,attr"` // This will hold user JID after auth.
	To               string `xml:"to,attr"`
	Version          string `xml:"version,attr"`
	DefaultNamespace string
	FeaturesState    features.State
	elements.ElementFactory
}

var Factory = elements.NewElementFactory()

func NewStream() *Stream {
	return &Stream{
		ElementFactory: Factory,
	}
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
