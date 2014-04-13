package roster

import (
	"encoding/xml"
	"log"

	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/stanzas/iq"
)

type RosterState struct {
	// here go methods to get / update the roster
}

type RosterItemElement struct {
	XMLName      xml.Name `xml:"item"`
	JID          string   `xml:"jid,attr,omitempty"`
	Name         string   `xml:"name,attr,omitempty"`
	Subscription string   `xml:"subscription,attr,omitempty"`
	Approved     bool     `xml:"approved,attr,omitempty"`
	Ask          string   `xml:"ask,attr,omitempty"`
	Groups       []string `xml:"group"`
}

type RosterElement struct {
	XMLName xml.Name `xml:"jabber:iq:roster query"`
	Ver     string   `xml:"ver,attr,omitempty"`
	Items   []RosterItemElement
}

func (self *RosterElement) Handle(request_id *iq.IQElement, stream *stream.Stream) error {
	// FIXME(dotdoom): 2014-04-03: auth check, state presence check, bind etc
	var state *RosterState
	stream.State.Get(&state)

	log.Printf("Roster request received")

	ri := RosterElement{}
	ri.Ver = "1.0"
	ri.Items = append(ri.Items, RosterItemElement{
		JID:          "test@localhost",
		Name:         "Tester",
		Subscription: "both",
	})
	ri.Items[0].Groups = append(ri.Items[0].Groups, "TestGroup")

	// TODO(dotdoom): 2014-04-03: might be easier to just use original IQ?
	response_iq := iq.NewIQElement()
	response_iq.Type = "result"
	response_iq.ID = request_id.ID
	response_iq.AddElement(&ri)
	stream.WriteElement(response_iq)

	return nil
}

func init() {
	iq.IQFactory.AddConstructor("jabber:iq:roster query", func() elements.Element {
		return &RosterElement{}
	})
}
