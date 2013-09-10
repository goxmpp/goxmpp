package stream

import "github.com/dotdoom/goxmpp/stream/decoder"

type Wrapper struct {
	rwStream      io.ReadWriter
	streamEncoder *xml.Encoder
	streamDecoder *xml.Decoder
	InnerDecoder  *decoder.InnerDecoder
	State         map[string]interface{}
}

func NewWrapper(rw io.ReadWriter) *StreamWrapper {
	xml_buffer := NewXMLBuffer()

	return &Wrapper{
		rwStream:      rw,
		streamEncoder: xml.NewEncoder(rw),
		streamDecoder: xml.NewDecoder(rw),
		InnerDecoder:  decoder.NewInnerDecoder(),
		State:         make(map[string]interface{}),
	}
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

func (sw *Wrapper) ReadXMLChunk(types map[[2]string](func(xml.StartElement) interface{})) (interface{}, error) {
	for {
		t, err := sw.Decoder.Token()
		if err != nil {
			return nil, err
		}
		if element, ok := t.(xml.StartElement); ok {
			// TODO(artem): handle </stream:stream> etc
			if tp, ok := types[[2]string{element.Name.Local, element.Name.Space}]; ok {
				value := tp(element)
				err = sw.Decoder.DecodeElement(value, &element)
				return value, err
			}
		}
	}
}

/*
func ReadStanza(d *xml.Decoder) (interface{}, error) {
	var element xml.StartElement
	for {
		t, err := d.Token()
		if err != nil { return nil, err }
		var ok bool
		element, ok = t.(xml.StartElement)
		var stanza interface{}
		if ok {
			switch element.Name.Local {
			case "message": stanza = new(Message)
			case "iq": stanza = new(IQ)
			case "presence": stanza = new(Presence)
			}
		}
		err = d.DecodeElement(stanza, &element)
		return stanza, err
	}
}
*/
