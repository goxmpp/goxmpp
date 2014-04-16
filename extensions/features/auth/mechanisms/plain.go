package mechanisms

import (
	"bytes"
	"encoding/base64"
	"errors"
	"log"

	"github.com/dotdoom/goxmpp/extensions/features/auth"
	"github.com/dotdoom/goxmpp/stream"
)

type PlainState struct {
	VerifyUserAndPassword func(string, string) bool
	RequireEncryption     bool
}

var usernamePasswordSeparator = []byte{0}

func init() {
	auth.AddMechanism("PLAIN", func(e *auth.AuthElement, stream *stream.Stream) error {
		b, err := base64.StdEncoding.DecodeString(e.Data)
		if err != nil {
			log.Println("Could not decode Base64 in Plain handler:", err)
			if err := stream.WriteElement(auth.NewFailute(IncorrectEncoding{})); err != nil {
				return err
			}
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

			if err := stream.WriteElement(&SuccessElement{}); err != nil {
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

	auth.MechanismsElement.AddElement(newMechanismElement("PLAIN"))
}
