package stream

import "github.com/dotdoom/goxmpp/stream/decoder"

type Wrapper struct {
	rwStream      io.ReadWriter
	StreamEncoder *xml.Encoder
	StreamDecoder *xml.Decoder
	InnerDecoder  *decoder.InnerDecoder
	State         map[string]interface{}
}

func NewWrapper(rw io.ReadWriter) *Wrapper {
	return &Wrapper{
		rwStream:      rw,
		StreamEncoder: xml.NewEncoder(rw),
		StreamDecoder: xml.NewDecoder(rw),
		InnerDecoder:  decoder.NewInnerDecoder(),
		State:         make(map[string]interface{}),
	}
}

func (self *Wrapper) SwapIOStream(rw io.ReadWriter) {
	self.rwStream = rw
	self.StreamEncoder = xml.NewEncoder(rw)
	self.StreamDecoder = xml.NewDecoder(rw)
}

func (sw *Wrapper) ReadStreamOpen() (*Stream, error) {
	stream := Stream{}
	for {
		t, err := sw.Decoder.Token()
		if err != nil {
			return nil, err
		}
		switch t := t.(type) {
		case xml.ProcInst:
			// Good.
		case xml.StartElement:
			if t.Name.Local == "stream" {
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
}

// TODO(artem): refactor
func (sw *Wrapper) WriteStreamOpen(stream *Stream, default_namespace string) (err error) {
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

	_, err = io.WriteString(sw.RW, data)
	return
}

func (sw *Wrapper) WriteFeatures() error {
	return nil
}
