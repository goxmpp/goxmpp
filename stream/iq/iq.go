package iq

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/stanza"

const (
	STREAD_NODE = "iq"
)

func init() {
	stream.HandlerRegistrator.Register(" "+STREAD_NODE, func() stream.Element {
		return &IQ{InnerXML: stream.InnerXML{Registrator: HandlerRegistrator}}
	})
}

var HandlerRegistrator = stream.NewElementGeneratorRegistrator()

/*type IQ struct {
	stream.Stanza
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
*/

type IQ struct {
	XMLName xml.Name `xml:"iq"`
	stanza.BaseStanza
	stream.InnerXML
}
