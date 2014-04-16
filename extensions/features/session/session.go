package session

import (
	"encoding/xml"
	"log"

	"github.com/dotdoom/goxmpp/extensions/features/bind"
	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
	"github.com/dotdoom/goxmpp/stream/elements/stanzas/iq"
)

type sessionFeatureElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-session session"`
}

type SessionElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-session session"`
}

type SessionState struct {
	Active bool
}

func (self *sessionFeatureElement) IsRequiredFor(stream *stream.Stream) bool {
	// According to RFC6121, which supersedes RFC3921, session MAY be supported.
	return false
}

func (self *sessionFeatureElement) CopyIfAvailable(stream *stream.Stream) elements.Element {
	var bind_state *bind.BindState
	err := stream.State.Get(&bind_state)
	if err == nil && bind_state.Resource != "" {
		var session_state *SessionState
		if err = stream.State.Get(&session_state); err != nil || !session_state.Active {
			return self
		}
	}
	return nil
}

func (self *SessionElement) Handle(request_id *iq.IQElement, stream *stream.Stream) error {
	// FIXME(dotdoom): 2014-04-04: auth check, state presence check, resource check required
	var state *SessionState
	if err := stream.State.Get(&state); err != nil {
		state = &SessionState{}
		stream.State.Push(state)
	}
	state.Active = true

	log.Printf("Session opened")

	// TODO(dotdoom): 2014-04-03: might be easier to just use original IQ?
	response_iq := iq.NewIQElement()
	response_iq.Type = "result"
	response_iq.ID = request_id.ID
	if err := stream.WriteElement(response_iq); err != nil {
		return err
	}

	return nil
}

func init() {
	iq.IQFactory.AddConstructor(func() elements.Element {
		return &SessionElement{}
	})
	features.Tree.AddElement(&sessionFeatureElement{})
}
