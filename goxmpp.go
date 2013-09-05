package goxmpp

import (
	"encoding/xml"
	"io"
)

type Stream struct {
	XMLName    xml.Name `xml:"http://etherx.jabber.org/streams stream"`
	ID         string   `xml:"id,attr,omitempty"`
	From       string   `xml:"from,attr,omitempty"`
	To         string   `xml:"to,attr,omitempty"`
	Version    string   `xml:"version,attr,omitempty"`
}

/*
 * Features
 */
type StreamFeature interface {
	ExposeTo(*StreamWrapper) StreamFeature
	AddSubfeature(StreamFeature) bool
	ClearSubfeatures()
}
type SimpleStreamFeature struct {
	Subfeatures []StreamFeature
}
func (ssf *SimpleStreamFeature) AddSubfeature(sf StreamFeature) bool {
	if sf != nil {
		ssf.Subfeatures = append(ssf.Subfeatures, sf)
		return true
	}
	return false
}
func (ssf *SimpleStreamFeature) ClearSubfeatures() {
	ssf.Subfeatures = nil
}
func (ssf *SimpleStreamFeature) ExposeSubfeaturesTo(sw *StreamWrapper, sf StreamFeature) StreamFeature {
	sf.ClearSubfeatures()
	for _, feature := range ssf.Subfeatures {
		sf.AddSubfeature(feature.ExposeTo(sw))
	}
	return sf
}

type StreamFeatures struct {
	XMLName xml.Name `xml:"stream:features"`
	SimpleStreamFeature
}
func (sf *StreamFeatures) ExposeTo(sw *StreamWrapper) StreamFeature {
	// TODO(artem): check availability
	return sf.ExposeSubfeaturesTo(sw, new(StreamFeatures))
}

type AuthMechanismsStreamFeature struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	SimpleStreamFeature
}
func (amsf *AuthMechanismsStreamFeature) ExposeTo(sw *StreamWrapper) StreamFeature {
	// TODO(artem): check availability
	return amsf.ExposeSubfeaturesTo(sw, new(AuthMechanismsStreamFeature))
}
type AuthMechanismStreamFeature struct {
	XMLName xml.Name `xml:"mechanism"`
	Name    string   `xml:",chardata"`
	SimpleStreamFeature
}
func (amsf *AuthMechanismStreamFeature) ExposeTo(*StreamWrapper) StreamFeature {
	// TODO(artem): check availability
	c_amsf := *amsf
	return &c_amsf
}

type CompressionMethodsStreamFeature struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl compression"`
	SimpleStreamFeature
}
func (cmsf *CompressionMethodsStreamFeature) ExposeTo(sw *StreamWrapper) StreamFeature {
	// TODO(artem): check availability
	return cmsf.ExposeSubfeaturesTo(sw, new(CompressionMethodsStreamFeature))
}
type CompressionMethodStreamFeature struct {
	XMLName xml.Name `xml:"method"`
	Name    string   `xml:",chardata"`
	SimpleStreamFeature
}
func (cmsf *CompressionMethodStreamFeature) ExposeTo(*StreamWrapper) StreamFeature {
	// TODO(artem): check availability
	c_cmsf := *cmsf
	return &c_cmsf
}

type StartTLSStreamFeature struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required bool     `xml:"required,omitempty"`
	SimpleStreamFeature
}
func (stsf *StartTLSStreamFeature) ExposeTo(*StreamWrapper) StreamFeature {
	// TODO(artem): check availability
	c_stsf := *stsf
	return &c_stsf
}

var GlobalStreamFeatures StreamFeatures
func RegisterGlobalStreamFeatures() {
	compressionMethods := CompressionMethodsStreamFeature{}
	compressionMethods.AddSubfeature(&CompressionMethodStreamFeature{Name: "zlib"})
	GlobalStreamFeatures.AddSubfeature(&compressionMethods)

	authMechanisms := AuthMechanismsStreamFeature{}
	authMechanisms.AddSubfeature(&AuthMechanismStreamFeature{Name: "PLAIN"})
	authMechanisms.AddSubfeature(&AuthMechanismStreamFeature{Name: "DIGEST-MD5"})
	GlobalStreamFeatures.AddSubfeature(&authMechanisms)

	startTLS := StartTLSStreamFeature{}
	startTLS.Required = true
	GlobalStreamFeatures.AddSubfeature(&startTLS)
}

/*
 * Stanzas
 */
type Error struct {
	Type   string `xml:"type,attr,omitempty"`
}

type Stanza struct {
	From   string `xml:"from,attr,omitempty"`
	To     string `xml:"to,attr,omitempty"`
	Type   string `xml:"type,attr,omitempty"`
	ID     string `xml:"id,attr,omitempty"`
	Error  Error
}

type Message struct {
	XMLName xml.Name `xml:"message"`
	Stanza
	Body    string   `xml:"body,omitempty"`
}

type VersionQuery struct {
	// http://xmpp.org/extensions/xep-0092.html
	XMLName xml.Name `xml:"jabber:iq:version query"`
	Name    string   `xml:"name,attr,omitempty"`
	Version string   `xml:"version,attr,omitempty"`
	OS      string   `xml:"os,attr,omitempty"`
}

type TimeQuery struct {
	// http://xmpp.org/extensions/xep-0202.html
	XMLName xml.Name `xml:"urn:xmpp:time time"`
	TZO     string   `xml:"tzo,omitempty"`
	UTC     string   `xml:"utc,omitempty"`
}

type DiscoInfoQuery struct {
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
}

type IQ struct {
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
}

type Presence struct {
	Stanza
	XMLName  xml.Name `xml:"presence"`
	Show     string   `xml:"show,omitempty"`
	Status   string   `xml:"status,omitempty"`
	Priority int      `xml:"priority,omitempty"`
}

type StreamWrapper struct {
	RW io.ReadWriter
	Encoder *xml.Encoder
	Decoder *xml.Decoder
}

func NewStreamWrapper(rw io.ReadWriter) *StreamWrapper {
	return &StreamWrapper{
		RW: rw,
		Encoder: xml.NewEncoder(rw),
		Decoder: xml.NewDecoder(rw) }
}

func (sw *StreamWrapper) ReadStreamOpen() (*Stream, error) {
	stream := Stream{}
	for {
		t, err := sw.Decoder.Token()
		if err != nil { return nil, err }
		switch t := t.(type) {
		case xml.ProcInst:
			// Good.
		case xml.StartElement:
			if t.Name.Local == "stream" {
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "to": stream.To = attr.Value
					case "from": stream.From = attr.Value
					case "version": stream.Version = attr.Value
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
		if err != nil { return nil, err }
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
