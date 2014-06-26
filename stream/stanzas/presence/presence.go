package presence

import (
	"encoding/xml"

	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/stanzas"
	"github.com/goxmpp/xtream"
)

var presenseXMLName = xml.Name{Local: "presence"}

func init() {
	xtream.NodeFactory.Add(func() xtream.Element {
		return NewPresenceElement()
	}, stream.StreamXMLName, presenseXMLName)
}

func NewPresenceElement() *PresenceElement {
	return &PresenceElement{InnerElements: xtream.NewElements(&presenseXMLName)}
}

type PresenceElement struct {
	XMLName xml.Name `xml:"presence"`
	Show    string   `xml:"show,omitempty"`
	Status  string   `xml:"status"`
	stanzas.Base
	xtream.InnerElements `xml:",any"`
}
