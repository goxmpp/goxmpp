package md5

import (
	"errors"

	"github.com/azhavnerchik/sasl/md5"
	"github.com/dotdoom/goxmpp/extensions/features/auth"
	"github.com/dotdoom/goxmpp/extensions/features/auth/mechanisms"
	"github.com/dotdoom/goxmpp/stream"
)

type DigestMD5Element string

func (md5 DigestMD5Element) IsAvailable(strm *stream.Stream) bool {
	var state *DigestMD5State
	if err := strm.State.Get(&state); err == nil {
		return true
	}
	return false
}

type DigestMD5State struct {
	Realm []string
	Host  string
}

type digestMD5Handler struct {
	state *DigestMD5State
	strm  *stream.Stream
	md5   *md5.MD5
}

func newDigestMD5Handler(state *DigestMD5State, strm *stream.Stream) (*digestMD5Handler, error) {
	md5 := md5.New()
	md5.SetRealm(state.Realm...)
	md5.SetQOP("auth")
	return &digestMD5Handler{md5: md5, state: state, strm: strm}, nil
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
	raw_resp_data, err := mechanisms.DecodeBase64(resp_el.Data, h.strm)
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
	auth.AddMechanism("DIGEST-MD5", func(e *auth.AuthElement, strm *stream.Stream) error {
		var state *DigestMD5State
		if err := strm.State.Get(&state); err != nil {
			return err
		}
		handler, err := newDigestMD5Handler(state, strm)
		if err != nil {
			return err
		}

		return handler.Handle()
	})

	auth.MechanismsElement.AddElement(mechanisms.NewMechanismElement(DigestMD5Element("DIGEST-MD5")))
}
