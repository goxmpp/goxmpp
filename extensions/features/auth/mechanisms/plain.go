package mechanisms

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"

	"github.com/dotdoom/goxmpp/extensions/features/auth"
	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
)

type SuccessElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
}

type PlainState struct {
	VerifyUserAndPassword func(string, string) bool
	RequireEncryption     bool
}

type PlainElement struct {
	XMLName xml.Name `xml:"mechanism"`
	Name    string   `xml:",chardata"`
}

func newPlainElement() *PlainElement {
	return &PlainElement{Name: "PLAIN"}
}

func (self *PlainElement) CopyIfAvailable(stream *stream.Stream) elements.Element {
	var plain_state *PlainState
	if err := stream.State.Get(&plain_state); err == nil {
		return self
	}
	return nil
}

var usernamePasswordSeparator = []byte{0}

func init() {
	auth.AddMechanism("PLAIN", func(e *auth.AuthElement, stream *stream.Stream) error {
		b, _ := base64.StdEncoding.DecodeString(e.Data)
		user_password := bytes.Split(b, usernamePasswordSeparator)

		var plain_state *PlainState
		if err := stream.State.Get(&plain_state); err != nil {
			return err
		}

		if plain_state.VerifyUserAndPassword(string(user_password[1]), string(user_password[2])) {
			var auth_state *auth.AuthState
			if err := stream.State.Get(&auth_state); err != nil {
				auth_state = &auth.AuthState{}
				stream.State.Push(auth_state)
			}
			auth_state.UserName = string(user_password[1])
			auth_state.Mechanism = "PLAIN"

			if err := stream.WriteElement(&SuccessElement{}); err != nil {
				return err
			}
			stream.ReOpen = true

			return nil
		} else {
			return errors.New("AUTH FAILED")
		}
	})
	auth.MechanismsElement.AddElement(newPlainElement())
}
