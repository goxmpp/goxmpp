package bind

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/extensions/features/auth"
	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
	"github.com/dotdoom/goxmpp/stream/elements/stanzas/iq"
)

type BindState struct {
	Resource string
}

type BindElement struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string   `xml:"resource"`
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
		return err != nil || state.Resource == ""
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

func init() {
	iq.IQFactory.AddConstructor("urn:ietf:params:xml:ns:xmpp-bind bind", func() elements.Element {
		return &BindElement{}
	})
	features.Tree.AddElement(&bindElement{})
}
