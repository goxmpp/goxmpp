package stream

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/decoder"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
	"io"
)

type Wrapper struct {
	rwStream       io.ReadWriter
	StreamEncoder  *xml.Encoder
	StreamDecoder  *xml.Decoder
	InnerDecoder   *decoder.InnerDecoder
	ElementFactory elements.Factory
	FeatureSet     *features.Features
	State          features.FeatureState
}

func NewWrapper(rw io.ReadWriter) *Wrapper {
	return &Wrapper{
		rwStream:       rw,
		StreamEncoder:  xml.NewEncoder(rw),
		StreamDecoder:  xml.NewDecoder(rw),
		InnerDecoder:   decoder.NewInnerDecoder(),
		State:          features.FeatureState{},
		ElementFactory: elements.GlobalFeaturesFactory,
		FeatureSet:     features.GlobalFeaturesList,
	}
}

func (self *Wrapper) SwapIOStream(rw io.ReadWriter) {
	self.rwStream = rw
	self.StreamEncoder = xml.NewEncoder(rw)
	self.StreamDecoder = xml.NewDecoder(rw)
}

func (self *Wrapper) FeaturesLoop() {
	for self.FeatureSet.IsRequiredFor(self.State) {
		self.StreamEncoder.Encode(self.FeatureSet.CopyIfAvailable(self.State))
		break
	}
}

func (sw *Wrapper) ReadStreamOpen() (*Stream, error) {
	for {
		t, err := sw.StreamDecoder.Token()
		if err != nil {
			return nil, err
		}
		switch t := t.(type) {
		case xml.ProcInst:
			// Good.
		case xml.StartElement:
			if t.Name.Local == "stream" {
				stream := Stream{}
				stream.XMLName = t.Name
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "to":
						stream.To = attr.Value
					case "from":
						stream.From = attr.Value
					case "version":
						stream.Version = attr.Value
					}
				}

				return &stream, nil
			}
		}
	}
	return nil, nil
}

// TODO(artem): refactor
func (sw *Wrapper) WriteStreamOpen(stream *Stream, default_namespace string) error {
	data := xml.Header

	data += "<stream:stream xmlns='" + default_namespace + "' xmlns:stream='" + stream.XMLName.Space + "'"
	if stream.ID != "" {
		data += " id='" + stream.ID + "'"
	}
	if stream.From != "" {
		data += " from='" + stream.From + "'"
	}
	if stream.To != "" {
		data += " to='" + stream.To + "'"
	}
	if stream.Version != "" {
		data += " version='" + stream.Version + "'"
	}
	data += ">"

	_, err := io.WriteString(sw.rwStream, data)
	return err
}
