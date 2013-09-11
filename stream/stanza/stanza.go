package stanza

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream"
)

type BaseStanza struct {
	XMLName xml.Name
	From    string `xml:"from,attr,omitempty"`
	To      string `xml:"to,attr,omitempty"`
	Type    string `xml:"type,attr,omitempty"`
	ID      string `xml:"id,attr,omitempty"`
	Lang    string `xml:"xml:lang,attr,omitempty"`
}

type StanzaWriter struct {
	BaseStanza
	stream.InnerElements
}

type ParsedStanza struct {
	BaseStanza
}
