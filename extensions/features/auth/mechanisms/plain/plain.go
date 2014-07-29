package plain

import (
	"bytes"
	"errors"

	"github.com/goxmpp/goxmpp/extensions/features/auth"
	"github.com/goxmpp/goxmpp/extensions/features/auth/mechanisms"
	"github.com/goxmpp/goxmpp/stream"
)

type PlainState struct {
	VerifyUserAndPassword func(string, string) bool
	RequireEncryption     bool
}

var usernamePasswordSeparator = []byte{0}

func init() {
	auth.AddMechanism("PLAIN",
		func(e *auth.AuthElement, stream *stream.Stream) error {
			var auth_state *auth.AuthState
			if err := stream.State.Get(&auth_state); err != nil {
				return err
			}

			b, err := auth.DecodeBase64(e.Data, stream)
			if err != nil {
				return err
			}
			user_password := bytes.Split(b, usernamePasswordSeparator)

			if pass := auth_state.GetPasswordByUserName(string(user_password[1])); pass == string(user_password[2]) {

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
