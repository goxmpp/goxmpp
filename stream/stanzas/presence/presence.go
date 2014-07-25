package presence

import (
	"encoding/xml"

	"github.com/goxmpp/goxmpp/stream/stanzas"
	"github.com/goxmpp/xtream"
)

func init() {
	xtream.NodeFactory.Add(func() xtream.Element {
		return NewPresenceElement()
	})
}

func NewPresenceElement() *PresenceElement {
	return &PresenceElement{InnerElements: xtream.NewElements()}
}

type PresenceElement struct {
	XMLName xml.Name `xml:"presence" parent:"stream:stream"`
	Show    string   `xml:"show,omitempty"`
	Status  string   `xml:"status"`
	stanzas.Base
	xtream.InnerElements `xml:",any"`
}
