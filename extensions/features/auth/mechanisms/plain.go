package mechanisms

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"

	"github.com/dotdoom/goxmpp/extensions/features/auth"
	"github.com/dotdoom/goxmpp/stream"
)

type SuccessElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
}

type PlainState struct {
	Callback          func(string, string) bool
	RequireEncryption bool
}

var usernamePasswordSeparator = []byte{0}

func init() {
	auth.AddMechanism("PLAIN", func(e *auth.AuthElement, stream *stream.Stream) error {
		b, _ := base64.StdEncoding.DecodeString(e.Data)
		user_password := bytes.Split(b, usernamePasswordSeparator)

		var plain_state *PlainState
		stream.State.Get(&plain_state)

		if plain_state.Callback(string(user_password[1]), string(user_password[2])) {
			var auth_state *auth.State
			if err := stream.State.Get(&auth_state); err != nil {
				auth_state = &auth.State{}
				stream.State.Push(auth_state)
			}
			auth_state.UserName = string(user_password[1])
			auth_state.Mechanism = "PLAIN"

			stream.WriteElement(&SuccessElement{})
			stream.Opened = false

			return nil
		} else {
			return errors.New("AUTH FAILED")
		}
	})
}
