package goxmpp

import (
	"encoding/xml"
	"io"
)

/*type Stream struct {
	XMLName    xml.Name `xml:"http://etherx.jabber.org/streams stream"`
	ID         string   `xml:"id,attr,omitempty"`
	From       string   `xml:"from,attr,omitempty"`
	To         string   `xml:"to,attr,omitempty"`
	Version    string   `xml:"version,attr,omitempty"`
}

/*
 * Features
*/
/*type StreamFeature interface {
	ExposeTo(*StreamWrapper) StreamFeature
	AddSubfeature(StreamFeature) bool
}
type SimpleStreamFeature struct {
	Subfeatures []StreamFeature
}
func (self *SimpleStreamFeature) AddSubfeature(sf StreamFeature) bool {
	if sf != nil {
		self.Subfeatures = append(self.Subfeatures, sf)
		return true
	}
	return false
}
func (self *SimpleStreamFeature) ExposeSubfeaturesTo(sw *StreamWrapper, sf StreamFeature) StreamFeature {
	for _, feature := range self.Subfeatures {
		sf.AddSubfeature(feature.ExposeTo(sw))
	}
	return sf
}

type StreamFeatures struct {
	XMLName xml.Name `xml:"stream:features"`
	SimpleStreamFeature
}
func (self *StreamFeatures) ExposeTo(sw *StreamWrapper) StreamFeature {
	if sw.State["session"] != nil { return nil }
	return self.ExposeSubfeaturesTo(sw, new(StreamFeatures))
}
*/

/*type AuthMechanismsStreamFeature struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	SimpleStreamFeature
}*/

func (self *AuthMechanismsStreamFeature) ExposeTo(sw *StreamWrapper) StreamFeature {
	if sw.State["authenticated"] != nil {
		return nil
	}
	return self.ExposeSubfeaturesTo(sw, new(AuthMechanismsStreamFeature))
}

/*type AuthMechanismStreamFeature struct {
	XMLName xml.Name `xml:"mechanism"`
	Name    string   `xml:",chardata"`
	SimpleStreamFeature
}*/

func (self *AuthMechanismStreamFeature) ExposeTo(*StreamWrapper) StreamFeature {
	c := *self
	return &c
}

/*type CompressionMethodsStreamFeature struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl compression"`
	SimpleStreamFeature
}*/

func (self *CompressionMethodsStreamFeature) ExposeTo(sw *StreamWrapper) StreamFeature {
	if sw.State["compressed"] != nil {
		return nil
	}
	return self.ExposeSubfeaturesTo(sw, new(CompressionMethodsStreamFeature))
}

/*type CompressionMethodStreamFeature struct {
	XMLName xml.Name `xml:"method"`
	Name    string   `xml:",chardata"`
	SimpleStreamFeature
}*/

func (self *CompressionMethodStreamFeature) ExposeTo(*StreamWrapper) StreamFeature {
	c := *self
	return &c
}

/*type StartTLSStreamFeature struct {
	XMLName     xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required    bool     `xml:"required,omitempty"`
	Certificate []byte
	SimpleStreamFeature
}*/

func (self *StartTLSStreamFeature) ExposeTo(sw *StreamWrapper) StreamFeature {
	if sw.State["tls"] != nil {
		return nil
	}
	if self.Certificate == nil {
		return nil
	}
	return self.ExposeSubfeaturesTo(sw, new(StartTLSStreamFeature))
}

/*type BindStreamFeature struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	SimpleStreamFeature
}

func (self *BindStreamFeature) ExposeTo(sw *StreamWrapper) StreamFeature {
	if sw.State["authenticated"] == nil {
		return nil
	}
	return self.ExposeSubfeaturesTo(sw, new(BindStreamFeature))
}*/

var GlobalStreamFeatures StreamFeatures

func RegisterGlobalStreamFeatures() {
	compressionMethods := CompressionMethodsStreamFeature{}
	compressionMethods.AddSubfeature(&CompressionMethodStreamFeature{Name: "zlib"})
	GlobalStreamFeatures.AddSubfeature(&compressionMethods)

	authMechanisms := AuthMechanismsStreamFeature{}
	authMechanisms.AddSubfeature(&AuthMechanismStreamFeature{Name: "PLAIN"})
	authMechanisms.AddSubfeature(&AuthMechanismStreamFeature{Name: "DIGEST-MD5"})
	GlobalStreamFeatures.AddSubfeature(&authMechanisms)

	GlobalStreamFeatures.AddSubfeature(new(SessionStreamFeature))
	GlobalStreamFeatures.AddSubfeature(new(BindStreamFeature))

	startTLS := StartTLSStreamFeature{}
	startTLS.Required = true
	GlobalStreamFeatures.AddSubfeature(&startTLS)
}

/*
 * Stanzas
 */
/*type Error struct {
	Type string `xml:"type,attr,omitempty"`
}*/

/*type Stanza struct {
	From  string `xml:"from,attr,omitempty"`
	To    string `xml:"to,attr,omitempty"`
	Type  string `xml:"type,attr,omitempty"`
	ID    string `xml:"id,attr,omitempty"`
	Error Error
}*/

/*type Message struct {
	XMLName xml.Name `xml:"message"`
	Stanza
	Body string `xml:"body,omitempty"`
}*/

/*type VersionQuery struct {
	// http://xmpp.org/extensions/xep-0092.html
	XMLName xml.Name `xml:"jabber:iq:version query"`
	Name    string   `xml:"name,attr,omitempty"`
	Version string   `xml:"version,attr,omitempty"`
	OS      string   `xml:"os,attr,omitempty"`
}*/

/*type TimeQuery struct {
	// http://xmpp.org/extensions/xep-0202.html
	XMLName xml.Name `xml:"urn:xmpp:time time"`
	TZO     string   `xml:"tzo,omitempty"`
	UTC     string   `xml:"utc,omitempty"`
}*/

/*type DiscoInfoQuery struct {
	// http://xmpp.org/extensions/xep-0030.html
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#info query"`
}

type DiscoItemsQuery struct {
	// http://xmpp.org/extensions/xep-0030.html
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#items query"`
}

type PingQuery struct {
	// http://xmpp.org/extensions/xep-0199.html
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}

type StatsQuery struct {
	// http://xmpp.org/extensions/xep-0039.html
	XMLName xml.Name `xml:"http://jabber.org/protocol/stats query"`
}

type LastQuery struct {
	// http://xmpp.org/extensions/xep-0012.html
	XMLName xml.Name `xml:"jabber:iq:last query"`
	Seconds int      `xml:"seconds,attr,omitempty"`
}

type PrivacyQuery struct {
	// http://xmpp.org/rfcs/rfc3921.html
	XMLName xml.Name `xml:"jabber:iq:privacy query"`
}*/

/*type IQ struct {
	Stanza
	XMLName         xml.Name `xml:"iq"`
	VersionQuery    VersionQuery
	TimeQuery       TimeQuery
	DiscoInfoQuery  DiscoInfoQuery
	DiscoItemsQuery DiscoItemsQuery
	PingQuery       PingQuery
	StatsQuery      StatsQuery
	LastQuery       LastQuery
	PrivacyQuery    PrivacyQuery
}*/

/*type Presence struct {
	Stanza
	XMLName  xml.Name `xml:"presence"`
	Show     string   `xml:"show,omitempty"`
	Status   string   `xml:"status,omitempty"`
	Priority int      `xml:"priority,omitempty"`
}*/

type XMLBuffer struct {
	buffer []byte
	pos    int
	size   int
}

func NewXMLBuffer() *XMLBuffer {
	// TODO: Figure out "good" starting size
	return &XMLBuffer{make([]byte, 100)}
}

func (self *XMLBuffer) PutXML(inner_xml []byte) {
	// Reallocate if buffer is too small
	if len(self.buffer) < len(inner_xml) {
		self.buffer = make([]byte, len(inner_xml))
	}

	self.size = copy(self.buffer, inner_xml)
	self.pos = 0
}

func (self *XMLBuffer) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	if self.pos >= self.size {
		return 0, io.EOF
	}

	n = copy(b, self.buffer[self.pos:self.size])
	self.pos += n
	return n, nil
}

type StreamWrapper struct {
	rwStream           io.ReadWriter
	streamEncoder      *xml.Encoder
	streamDecoder      *xml.Decoder
	InnerDecoder       *xml.Decoder
	InnerDecoderBuffer XMLBuffer
	State              map[string]interface{}
}

func NewStreamWrapper(rw io.ReadWriter) *StreamWrapper {
	xml_buffer := NewXMLBuffer()

	return &StreamWrapper{
		rwStream:           rw,
		streamEncoder:      xml.NewEncoder(rw),
		streamDecoder:      xml.NewDecoder(rw),
		InnerDecoder:       xml.NewDecoder(xml_buffer),
		InnerDecoderBuffer: xml_buffer,
		State:              make(map[string]interface{}),
	}
}

func (sw *StreamWrapper) ReadStreamOpen() (*Stream, error) {
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
func (sw *StreamWrapper) WriteStreamOpen(stream *Stream, default_namespace string) (err error) {
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

func (sw *StreamWrapper) WriteFeatures() error {
	return nil
}

func (sw *StreamWrapper) ReadXMLChunk(types map[[2]string](func(xml.StartElement) interface{})) (interface{}, error) {
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
