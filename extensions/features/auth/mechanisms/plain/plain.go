package plain

import (
	"bytes"
	"errors"

	"github.com/goxmpp/goxmpp/extensions/features/auth"
	"github.com/goxmpp/goxmpp/extensions/features/auth/mechanisms"
	"github.com/goxmpp/goxmpp/stream"
)

type PlainElement string

type PlainState struct {
	VerifyUserAndPassword func(string, string) bool
	RequireEncryption     bool
}

var usernamePasswordSeparator = []byte{0}

func init() {
	auth.AddMechanism(mechanisms.NewMechanismElement(PlainElement("PLAIN")),
		func(e *auth.AuthElement, stream *stream.Stream) error {
			b, err := auth.DecodeBase64(e.Data, stream)
			if err != nil {
				return err
			}
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

				if err := stream.WriteElement(mechanisms.SuccessElement{}); err != nil {
					return err
				}

				auth_state.UserName = string(user_password[1])
				auth_state.Mechanism = "PLAIN"
				stream.ReOpen = true

				return nil
			} else {
				return errors.New("AUTH FAILED")
			}
		})
}
