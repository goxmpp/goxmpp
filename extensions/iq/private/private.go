package roster

import (
	"encoding/xml"
	"log"

	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/stanzas/iq"
)

type PrivateElement struct {
	XMLName xml.Name `xml:"jabber:iq:private query"`
	*elements.InnerElements
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

var PrivateFactory = elements.NewFactory()

func NewPrivateElement() *PrivateElement {
	return &PrivateElement{InnerElements: elements.NewInnerElements(PrivateFactory)}
}

func init() {
	iq.IQFactory.AddConstructor(func() elements.Element {
		return NewPrivateElement()
	})
}
