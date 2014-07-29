package session

import (
	"encoding/xml"
	"log"

	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/features"
	"github.com/goxmpp/goxmpp/stream/stanzas/iq"
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

func (self *SessionElement) Handle(strm *stream.Stream, opts features.Options) error {
	request_id := opts.(*iq.IQElement)
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
			return features.NewFeature("session", &sessionFeatureElement{}, false, nil)
		},
		Name:   xml.Name{Local: "session", Space: "urn:ietf:params:xml:ns:xmpp-session"},
		Parent: iq.IQXMLName,
		Wants:  []string{"bind"},
	})
}
