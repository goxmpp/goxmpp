package roster

import (
	"encoding/xml"
	"log"

	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/stanzas/iq"
	"github.com/goxmpp/xtream"
)

var privateXMLName = xml.Name{Local: "query", Space: "jabber:iq:private"}

type PrivateElement struct {
	XMLName              xml.Name `xml:"jabber:iq:private query" parent:"iq"`
	xtream.InnerElements `xml:",any"`
}

func (self *PrivateElement) Handle(request_id *iq.IQElement, stream *stream.Stream) error {
	log.Printf("Private storage request received")

	response_iq := iq.NewIQElement()
	response_iq.Type = "error"
	response_iq.ID = request_id.ID
	if err := stream.WriteElement(response_iq); err != nil {
		return err
	}

	return nil
}

func NewPrivateElement() *PrivateElement {
	return &PrivateElement{InnerElements: xtream.NewElements()}
}

func init() {
	xtream.NodeFactory.Add(func() xtream.Element {
		return NewPrivateElement()
	})
}
