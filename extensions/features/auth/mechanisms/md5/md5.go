package md5

import (
	"errors"

	"github.com/goxmpp/goxmpp/extensions/features/auth"
	"github.com/goxmpp/goxmpp/extensions/features/auth/mechanisms"
	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/sasl/digest"
)

type digestMD5Handler struct {
	strm *stream.Stream
	md5  *digest.Server
}

func newDigestMD5Handler(strm *stream.Stream) (*digestMD5Handler, error) {
	md5 := digest.NewServer(&digest.Options{Realms: []string{"test"}, QOPs: []string{"auth"}})
	return &digestMD5Handler{md5: md5, strm: strm}, nil
}

func (h *digestMD5Handler) Handle() error {
	var auth_state *auth.AuthState
	if err := h.strm.State.Get(&auth_state); err != nil {
		auth_state = &auth.AuthState{}
		h.strm.State.Push(auth_state)
	}

	if err := h.strm.WriteElement(mechanisms.NewChallengeElement(h.md5.Challenge())); err != nil {
		return err
	}

	// Receive a response with encoded MD5
	resp_el, err := mechanisms.ReadResponse(h.strm)
	if err != nil {
		return err
	}

	// Check MD5
	raw_resp_data, err := auth.DecodeBase64(resp_el.Data, h.strm)
	if err != nil {
		return err
	}

	if err := h.md5.ParseResponse(raw_resp_data); err != nil {
		return err
	}
	password := auth_state.GetPasswordByUserName(h.md5.UserName())
	if err := h.md5.Validate(password); err != nil {
		return err
	}

	// Send response
	if err := h.strm.WriteElement(mechanisms.NewChallengeElement(h.md5.Final())); err != nil {
		return err
	}

	rsp, err := mechanisms.ReadResponse(h.strm)
	if err != nil {
		return err
	}
	if rsp.Data != "" {
		return errors.New("Wrong response, expected empty response")
	}

	if err := h.strm.WriteElement(mechanisms.SuccessElement{}); err != nil {
		return err
	}

	auth_state.UserName = h.md5.AuthID()

	h.strm.ReOpen = true

	return nil
}

func init() {
	auth.AddMechanism("DIGEST-MD5",
		func(e *auth.AuthElement, strm *stream.Stream) error {
			handler, err := newDigestMD5Handler(strm)
			if err != nil {
				return err
			}

			return handler.Handle()
		})
}
