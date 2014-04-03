package bind

import (
	"encoding/xml"
	"log"

	"github.com/dotdoom/goxmpp/extensions/features/auth"
	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
	"github.com/dotdoom/goxmpp/stream/elements/stanzas/iq"
)

type BindState struct {
	Resource       string
	VerifyResource func(string) bool
}

type BindElement struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string   `xml:"resource,omitempty"`
	JID      string   `xml:"jid,omitempty"`
}

type bindElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
}

func (self *bindElement) IsRequiredFor(stream *stream.Stream) bool {
	var auth_state *auth.AuthState
	err := stream.State.Get(&auth_state)
	if err == nil && auth_state.UserName != "" {
		var state *BindState
		err := stream.State.Get(&state)
		// Bind is required by XMPP standard, so it has to be always present.
		// However we cheat here a little and allow to skip this.
		// FIXME(dotdoom): 2014-04-03: what will happen to JID in stream.{From,To} ?
		return err == nil && state.Resource == ""
	} else {
		return false
	}
}

func (self *bindElement) CopyIfAvailable(stream *stream.Stream) elements.Element {
	if self.IsRequiredFor(stream) {
		return &bindElement{}
	}
	return nil
}

func (self *BindElement) Handle(request_id *iq.IQElement, stream *stream.Stream) error {
	// FIXME(dotdoom): 2014-04-03: auth check, state presence check, resource check required
	var state *BindState
	stream.State.Get(&state)
	if state.VerifyResource(self.Resource) {
		state.Resource = self.Resource
	} else {
		// TODO(dotdoom): 2014-04-03
	}

	var authState *auth.AuthState
	stream.State.Get(&authState)

	stream.To = authState.UserName + "@" + stream.From + "/" + state.Resource
	log.Printf("Bound to JID: %#v", stream.To)

	// TODO(dotdoom): 2014-04-03: might be easier to just use original IQ?
	response_iq := iq.NewIQElement()
	response_iq.Type = "result"
	response_iq.ID = request_id.ID
	response_iq.AddElement(&BindElement{JID: stream.To})
	stream.WriteElement(response_iq)

	return nil
}

func init() {
	iq.IQFactory.AddConstructor("urn:ietf:params:xml:ns:xmpp-bind bind", func() elements.Element {
		return &BindElement{}
	})
	features.Tree.AddElement(&bindElement{})
}
