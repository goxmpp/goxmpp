package bind

import (
	"encoding/xml"
	"log"

	"github.com/goxmpp/goxmpp/extensions/features/auth"
	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/features"
	"github.com/goxmpp/goxmpp/stream/stanzas/iq"
)

type BindState struct {
	Resource       string
	VerifyResource func(string) bool
}

type BindElement struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind" parent:"iq"`
	Resource string   `xml:"resource,omitempty"`
	JID      string   `xml:"jid,omitempty"`
}

type bindElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
}

func (be *bindElement) NewHandler() features.FeatureHandler {
	return &BindElement{}
}

func (self *BindElement) Handle(strm *stream.Stream, opts features.Options) error {
	request_id := opts.(*iq.IQElement)

	// FIXME(goxmpp): 2014-04-03: auth check, state presence check, resource check required
	var state *BindState
	strm.State.Get(&state)
	if state.VerifyResource(self.Resource) {
		state.Resource = self.Resource
	} else {
		// TODO(goxmpp): 2014-04-03
	}

	var authState *auth.AuthState
	strm.State.Get(&authState)

	strm.To = authState.UserName + "@" + strm.From + "/" + state.Resource
	log.Printf("Bound to JID: %#v", strm.To)

	// TODO(goxmpp): 2014-04-03: might be easier to just use original IQ?
	response_iq := iq.NewIQElement()
	response_iq.Type = "result"
	response_iq.ID = request_id.ID
	response_iq.AddElement(&BindElement{JID: strm.To})
	if err := strm.WriteElement(response_iq); err != nil {
		return err
	}

	return nil
}

func init() {
	features.FeatureFactory.Add("bind", &features.FeatureFactoryElement{
		Constructor: func(opts features.Options) *features.Feature {
			return features.NewFeature("bind", &bindElement{}, true)
		},
		Name:   xml.Name{Local: "bind", Space: "urn:ietf:params:xml:ns:xmpp-bind"},
		Parent: iq.IQXMLName,
		Wants:  []string{"auth"},
	})
}
