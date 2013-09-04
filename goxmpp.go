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
}

type CompositeStreamFeature interface {
	AddSubfeature(StreamFeature) bool
	ClearSubfeatures()
	StreamFeature
}
type StreamFeaturesHolder struct {
	Subfeatures []StreamFeature
}
func (sfh *StreamFeaturesHolder) AddSubfeature(f StreamFeature) bool {
	if f != nil {
		sfh.Subfeatures = append(sfh.Subfeatures, f)
		return true
	}
	return false
}
func (sfh *StreamFeaturesHolder) ClearSubfeatures() {
	sfh.Subfeatures = sfh.Subfeatures[:0]
}
func (sfh *StreamFeaturesHolder) ExposeSubfeatures(sw *StreamWrapper, csf CompositeStreamFeature) {
	csf.ClearSubfeatures()
	for _, feature := range sfh.Subfeatures {
		csf.AddSubfeature(feature.ExposeTo(sw))
	}
}

type XMPPStreamFeatures struct {
	XMLName     xml.Name `xml:"stream:features"`
	StreamFeaturesHolder
}
func (xsf *XMPPStreamFeatures) ExposeTo(sw *StreamWrapper) StreamFeature {
	c_xsf := *xsf
	xsf.ExposeSubfeatures(sw, &c_xsf)
	return &c_xsf
}

type AuthMechanismsStreamFeature struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	StreamFeaturesHolder
}
func (amsf *AuthMechanismsStreamFeature) ExposeTo(sw *StreamWrapper) StreamFeature {
	if false {
		c_amsf := *amsf
		amsf.ExposeSubfeatures(sw, &c_amsf)
		return &c_amsf
	}
	return nil
}
type AuthMechanismStreamFeature struct {
	XMLName xml.Name `xml:"mechanism"`
	Name    string   `xml:",chardata"`
	StreamFeaturesHolder
}
func (amsf *AuthMechanismStreamFeature) ExposeTo(*StreamWrapper) StreamFeature {
	c_amsf := *amsf
	return &c_amsf
}

type CompressionMethodsStreamFeature struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl compression"`
	StreamFeaturesHolder
}
func (cmsf *CompressionMethodsStreamFeature) ExposeTo(sw *StreamWrapper) StreamFeature {
	c_cmsf := *cmsf
	cmsf.ExposeSubfeatures(sw, &c_cmsf)
	return &c_cmsf
}
type CompressionMethodStreamFeature struct {
	XMLName xml.Name `xml:"method"`
	Name    string   `xml:",chardata"`
	StreamFeaturesHolder
}
func (cmsf *CompressionMethodStreamFeature) ExposeTo(*StreamWrapper) StreamFeature {
	c_cmsf := *cmsf
	return &c_cmsf
}

type StartTLSStreamFeature struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required bool     `xml:"required,omitempty"`
	StreamFeaturesHolder
}
func (stsf *StartTLSStreamFeature) ExposeTo(*StreamWrapper) StreamFeature {
	c_stsf := *stsf
	return &c_stsf
}

var StreamFeatures XMPPStreamFeatures
func RegisterFeatures() {
	compressionMethods := CompressionMethodsStreamFeature{}
	compressionMethods.AddSubfeature(&CompressionMethodStreamFeature{Name: "zlib"})
	StreamFeatures.AddSubfeature(&compressionMethods)

	authMechanisms := AuthMechanismsStreamFeature{}
	authMechanisms.AddSubfeature(&AuthMechanismStreamFeature{Name: "PLAIN"})
	authMechanisms.AddSubfeature(&AuthMechanismStreamFeature{Name: "DIGEST-MD5"})
	StreamFeatures.AddSubfeature(&authMechanisms)

	startTLS := StartTLSStreamFeature{}
	StreamFeatures.AddSubfeature(&startTLS)
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
