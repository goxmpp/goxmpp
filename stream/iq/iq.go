package iq

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"

type IQ struct {
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
