package session

import (
	"encoding/xml"
	"log"

	"github.com/goxmpp/goxmpp/extensions/features/bind"
	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/features"
	"github.com/goxmpp/goxmpp/stream/stanzas/iq"
	"github.com/goxmpp/xtream"
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

func (self *sessionFeatureElement) CopyIfAvailable(stream *stream.Stream) xtream.Element {
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

func (self *SessionElement) Handle(fc features.FeatureContainable, opts features.Options) error {
	request_id := opts.(*iq.IQElement)
	strm := fc.(*stream.Stream)
	// FIXME(goxmpp): 2014-04-04: auth check, state presence check, resource check required
	var state *SessionState
	if err := strm.State.Get(&state); err != nil {
		state = &SessionState{}
		strm.State.Push(state)
	}
	state.Active = true

	log.Printf("Session opened")

	// TODO(goxmpp): 2014-04-03: might be easier to just use original IQ?
	response_iq := iq.NewIQElement()
	response_iq.Type = "result"
	response_iq.ID = request_id.ID
	if err := strm.WriteElement(response_iq); err != nil {
		return err
	}

	return nil
}

func (s *sessionFeatureElement) NewHandler() features.FeatureHandler {
	return &SessionElement{}
}

func init() {
	features.FeatureFactory.Add("session", &features.FeatureFactoryElement{
		Constructor: func(opts features.Options) *features.Feature {
			return features.NewFeature("session", &sessionFeatureElement{}, false)
		},
		Name:   xml.Name{Local: "session", Space: "urn:ietf:params:xml:ns:xmpp-session"},
		Parent: iq.IQXMLName,
	})
}
