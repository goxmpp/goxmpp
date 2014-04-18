package md5

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/dotdoom/goxmpp/extensions/features/auth"
	"github.com/dotdoom/goxmpp/extensions/features/auth/mechanisms"
	"github.com/dotdoom/goxmpp/stream"
)

const A2_AUTH_SUFFIX = "00000000000000000000000000000000"

type DigestMD5Element string

func (md5 DigestMD5Element) IsAvailable(strm *stream.Stream) bool {
	var state *DigestMD5State
	if err := strm.State.Get(&state); err == nil {
		return true
	}
	return false
}

type Challenge struct {
	Realm     []string
	Nonce     string
	QOP       string
	Charset   string
	Algorithm string
}

func NewChallenge(realm []string) (*Challenge, error) {
	nonce := make([]byte, 14)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return &Challenge{
		Realm:     realm,
		Nonce:     fmt.Sprintf("%x", nonce),
		QOP:       "auth",
		Algorithm: "md5-sess",
	}, nil
}

func (c *Challenge) String() string {
	str := []string{fmt.Sprintf("nonce=\"%s\"", c.Nonce), fmt.Sprintf("algorithm=\"%s\"", c.Algorithm)}
	for _, realm := range c.Realm {
		str = append(str, fmt.Sprintf("realm=\"%s\"", realm))
	}
	if c.QOP != "" {
		str = append(str, fmt.Sprintf("qop=\"%s\"", c.QOP))
	}
	if c.Charset != "" {
		str = append(str, fmt.Sprintf("charset=\"%s\"", c.QOP))
	}

	return strings.Join(str, ",")
}

type Response struct {
	UserName  string
	Realm     string
	Nonce     string
	CNonce    string
	NC        string
	ServType  string
	Host      string
	DigestURI string
	Response  string
	Charset   string
	AuthId    string
}

func NewMD5ResponseFromData(data []byte) *Response {
	resp := &Response{}

	for _, param := range bytes.Split(data, []byte(",")) {
		key_val := bytes.SplitN(param, []byte("="), 2)

		val := string(bytes.Trim(key_val[1], "\""))
		switch string(key_val[0]) {
		case "username":
			resp.UserName = val
		case "realm":
			resp.Realm = val
		case "nonce":
			resp.Nonce = val
		case "cnonce":
			resp.CNonce = val
		case "nc":
			resp.NC = val
		case "serv-type":
			resp.ServType = val
		case "host":
			resp.Host = val
		case "digest-uri":
			resp.DigestURI = val
		case "response":
			resp.Response = val
		case "charset":
			resp.Charset = val
		case "authzid":
			resp.AuthId = val
		}
	}

	return resp
}

func (r *Response) Validate(c *Challenge, state *DigestMD5State) error {
	// TODO check authid
	if r.Nonce != c.Nonce {
		return errors.New("Wrong nonce replied")
	}

	if r.Host != state.Host {
		return errors.New("Wrong host replied")
	}

	if len(c.Realm) > 0 {
		for _, realm := range c.Realm {
			if r.Realm == realm {
				return nil
			}
		}

		return errors.New("Wrong realm received from client")
	}

	return nil
}

type DigestMD5State struct {
	ValidateMD5 func(*Challenge, *Response) bool
	Realm       []string
	Host        string
}

type digestMD5Handler struct {
	state     *DigestMD5State
	challenge *Challenge
	strm      *stream.Stream
}

func newDigestMD5Handler(state *DigestMD5State, strm *stream.Stream) (*digestMD5Handler, error) {

	challenge, err := NewChallenge(state.Realm)
	if err != nil {
		log.Println("Could not create a challenge")
		return nil, err
	}

	return &digestMD5Handler{challenge: challenge, state: state, strm: strm}, nil
}

func (h *digestMD5Handler) Handle() error {
	if err := h.strm.WriteElement(mechanisms.NewChallengeElement(h.challenge.String())); err != nil {
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

	log.Println("Sent challenge", h.challenge.String())
	log.Println("Received response", string(raw_resp_data))

	resp := NewMD5ResponseFromData(raw_resp_data)
	log.Printf("Challenge object %#v", h.challenge)
	log.Printf("Response object %#v", resp)

	if err := resp.Validate(h.challenge, h.state); err != nil {
		return err
	}
	if !h.state.ValidateMD5(h.challenge, resp) {
		return errors.New("AUTH FAILED")
	}

	// Send response
	if err := h.strm.WriteElement(mechanisms.NewChallengeElement("rspauth")); err != nil {
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

	var auth_state *auth.AuthState
	if err := h.strm.State.Get(&auth_state); err != nil {
		auth_state = &auth.AuthState{}
		h.strm.State.Push(auth_state)
	}

	auth_state.UserName = resp.UserName

	h.strm.ReOpen = true

	return nil
}

func (r *Response) GenerateHash(c *Challenge, password string) string {
	x := md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", r.UserName, r.Realm, password)))

	start_str := fmt.Sprintf("%s:%s:%s", x, c.Nonce, r.CNonce)
	if r.AuthId != "" {
		start_str = fmt.Sprintf("%s:%s", start_str, r.AuthId)
	}
	start := md5.Sum([]byte(start_str))

	end_str := fmt.Sprintf("AUTHENTICATE:%s", r.DigestURI)
	if c.QOP == "auth-int" || c.QOP == "auth-conf" {
		end_str += fmt.Sprintf("%s:%s", end_str, A2_AUTH_SUFFIX)
	}
	end := md5.Sum([]byte(end_str))

	hash_str := fmt.Sprintf("%x:%s:%s:%s:%s:%x", start, c.Nonce, r.NC, r.CNonce, c.QOP, end)
	return fmt.Sprintf("%x", md5.Sum([]byte(hash_str)))
}

func init() {
	auth.AddMechanism("DIGEST-MD5", func(e *auth.AuthElement, strm *stream.Stream) error {
		var state *DigestMD5State
		if err := strm.State.Get(&state); err != nil {
			return nil
		}
		handler, err := newDigestMD5Handler(state, strm)
		if err != nil {
			return err
		}

		return handler.Handle()
	})

	auth.MechanismsElement.AddElement(mechanisms.NewMechanismElement(DigestMD5Element("DIGEST-MD5")))
}
