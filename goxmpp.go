package goxmpp

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream"
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

