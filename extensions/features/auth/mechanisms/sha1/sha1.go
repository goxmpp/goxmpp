package sha1

import (
	"crypto/sha1"
	"log"

	"github.com/goxmpp/goxmpp/extensions/features/auth"
	"github.com/goxmpp/goxmpp/extensions/features/auth/mechanisms"
	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/sasl/scram"
)

const MIN_ITERS = 4096

type SHAElement string

func (sha SHAElement) IsAvailable(strm *stream.Stream) bool {
	var state *SHAState
	if err := strm.State.Get(&state); err == nil && !state.Authenticated {
		return true
	}
	return false
}

type SHAState struct {
	Authenticated bool
}

type shaHandler struct {
	scram     *scram.Server
	strm      *stream.Stream
	authState *auth.AuthState
	shaState  *SHAState
}

func newSHAHandler(strm *stream.Stream, scram *scram.Server, astate *auth.AuthState, sstate *SHAState) *shaHandler {
	return &shaHandler{strm: strm, scram: scram, authState: astate, shaState: sstate}
}

func (h *shaHandler) Handle() error {
	if err := h.strm.WriteElement(mechanisms.NewChallengeElement(h.scram.First())); err != nil {
		return err
	}

	// Receive a response with encoded MD5
	resp_el, err := mechanisms.ReadResponse(h.strm)
	if err != nil {
		return err
	}

	// Check SHA
	raw_resp_data, err := mechanisms.DecodeBase64(resp_el.Data, h.strm)
	if err != nil {
		return err
	}

	if err := h.scram.CheckClientFinal(raw_resp_data); err != nil {
		return err
	}

	// Send response
	if err := h.strm.WriteElement(mechanisms.NewSuccessElement(h.scram.Final())); err != nil {
		log.Println("Could not write signature")
		return err
	}

	h.authState.UserName = h.scram.UserName()

	h.shaState.Authenticated = true

	h.strm.ReOpen = true

	return nil
}

func init() {
	auth.AddMechanism("SCRAM-SHA-1", func(e *auth.AuthElement, strm *stream.Stream) error {
		var state *SHAState
		if err := strm.State.Get(&state); err != nil {
			log.Println("SCRAM-SHA-1 is not available but tried to be used")
			return err
		}

		var auth_state *auth.AuthState
		if err := strm.State.Get(&auth_state); err != nil {
			log.Println("SHAM-SHA-1 AuthState is not set can't get auth data")
			return err
		}

		auth_data, err := mechanisms.DecodeBase64(e.Data, strm)
		if err != nil {
			return err
		}

		scram := scram.NewServer(sha1.New, nil)
		if err := scram.ParseClientFirst(auth_data); err != nil {
			return err
		}
		scram.SaltPassword([]byte(auth_state.GetPasswordByUserName(scram.UserName())))

		handler := newSHAHandler(strm, scram, auth_state, state)

		return handler.Handle()
	})

	auth.MechanismsElement.AddElement(mechanisms.NewMechanismElement(SHAElement("SCRAM-SHA-1")))
}
