package sha1

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	mrand "math/rand"
	"strings"
	"time"

	"github.com/dotdoom/goxmpp/extensions/features/auth"
	"github.com/dotdoom/goxmpp/extensions/features/auth/mechanisms"
	"github.com/dotdoom/goxmpp/stream"
)

type SHAElement string

func (sha SHAElement) IsAvailable(strm *stream.Stream) bool {
	var state *SHAState
	if err := strm.State.Get(&state); err == nil {
		return true
	}
	return false
}

type ClientAuth struct {
	Binding  string
	Nonce    []byte
	UserName string
}

func NewClientAuth(username, binding string) (*ClientAuth, error) {
	nonce := make([]byte, 50)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return &ClientAuth{
		Nonce:    nonce,
		UserName: username,
		Binding:  binding,
	}, nil
}

func NewClientAuthFromData(client_auth string) (*ClientAuth, error) {
	ca := &ClientAuth{}

	log.Println("Client auth string", client_auth)

	if err := CheckSCRAMString(client_auth); err != nil {
		return nil, err
	}

	tokens := strings.Split(client_auth, ",")
	binding := []string{}
	for _, tok := range tokens {
		if len(tok) <= 1 || tok[1] != '=' {
			binding = append(binding, tok)
		} else {
			kv := strings.SplitN(tok, "=", 2)
			switch kv[0] {
			case "n":
				ca.UserName = kv[1]
			case "r":
				log.Println("Nonce from auth string", kv[1])
				nonce, err := base64.StdEncoding.DecodeString(kv[1])
				if err != nil {
					return nil, err
				}
				ca.Nonce = nonce
			}
		}
	}
	ca.Binding = strings.Join(binding, ",")

	log.Printf("%#v", ca)
	return ca, nil
}

func (ca *ClientAuth) String() string {
	return fmt.Sprintf("%s,n=%s,r=%s", ca.Binding, ca.UserName, base64.StdEncoding.EncodeToString(ca.Nonce))
}

func CheckSCRAMString(scram string) error {
	if scram[0] != 'n' && scram[0] != 'c' && scram[0] != 'p' {
		return errors.New("Wrong SCRAM message")
	}
	return nil
}

type Challenge struct {
	Iterations int    // Iterations count
	Nonce      []byte // Client's + server's
	Salt       []byte // Encoded in Base64
}

func NewChallenge(base_nonce []byte) (*Challenge, error) {
	nonce := make([]byte, 50)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	salt := make([]byte, 30)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	mrand.Seed(time.Now().UnixNano())

	return &Challenge{
		Nonce:      append(base_nonce, nonce...),
		Salt:       salt,
		Iterations: 1024 + mrand.Intn(9000-1024),
	}, nil
}

func (c *Challenge) String() string {
	return strings.Join([]string{
		fmt.Sprintf("r=%s", base64.StdEncoding.EncodeToString(c.Nonce)),
		fmt.Sprintf("s=%s", base64.StdEncoding.EncodeToString(c.Salt)),
		fmt.Sprintf("i=%d", c.Iterations),
	}, ",")
}

type Response struct {
	Nonce      []byte
	Proof      []byte
	Binding    string // Base64 encoded
	server_sig []byte
}

func NewSHAResponseFromData(data []byte) (*Response, error) {
	if err := CheckSCRAMString(string(data)); err != nil {
		return nil, err
	}

	resp := &Response{}
	for _, token := range bytes.Split(data, []byte(",")) {
		kv := bytes.SplitN(token, []byte("="), 2)
		val := string(kv[1])
		switch string(kv[0]) {
		case "c":
			resp.Binding = val
		case "r":
			log.Println("Raw nonce", val)
			nonce, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				log.Println("Could not decode nonce")
				return nil, err
			}
			resp.Nonce = nonce
		case "p":
			proof, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				log.Println("Could not decode proof")
				return nil, err
			}
			resp.Proof = proof
		}
	}

	return resp, nil
}

func (r *Response) Validate(c *Challenge, state *SHAState) error {
	// TODO Added validation logic here
	return nil
}

type SHAState struct {
	Validate func(*Challenge, *Response) bool
}

type shaHandler struct {
	challenge   *Challenge
	client_auth *ClientAuth
	strm        *stream.Stream
}

func newSHAHandler(strm *stream.Stream, client_auth *ClientAuth) (*shaHandler, error) {
	challenge, err := NewChallenge(client_auth.Nonce)
	if err != nil {
		log.Println("Could not create a challenge")
		return nil, err
	}

	return &shaHandler{challenge: challenge, strm: strm, client_auth: client_auth}, nil
}

func (h *shaHandler) Handle() error {
	log.Println("Sending challenge", h.challenge)
	if err := h.strm.WriteElement(mechanisms.NewChallengeElement(h.challenge.String())); err != nil {
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

	log.Println("Sent challenge", h.challenge.String())
	log.Println("Received response", string(raw_resp_data))

	resp, err := NewSHAResponseFromData(raw_resp_data)
	if err != nil {
		return err
	}
	log.Printf("Challenge object %#v", h.challenge)
	log.Printf("Response object %#v", resp)

	//if err := resp.Validate(h.challenge); err != nil {
	//	return err
	//}
	var auth_state *auth.AuthState
	if err := h.strm.State.Get(&auth_state); err != nil {
		log.Println("Could not get auth state")
		return err
	}
	proof := resp.GenerateProof(h.challenge, h.client_auth, auth_state.GetPasswordByUserName(h.client_auth.UserName))
	log.Println("Expected proof", proof)
	log.Println("     Got proof", resp.Proof)
	if fmt.Sprintf("%x", proof) != fmt.Sprintf("%x", resp.Proof) {
		return errors.New("AUTH FAILED")
	}

	// Send response
	if err := h.strm.WriteElement(mechanisms.NewSuccessElement(GetServerSignatureMessage(resp.server_sig))); err != nil {
		return err
	}

	auth_state = &auth.AuthState{}
	h.strm.State.Push(auth_state)

	auth_state.UserName = h.client_auth.UserName

	h.strm.ReOpen = true

	return nil
}

func (c *Challenge) SaltPassword(password string) []byte {
	mac := hmac.New(sha1.New, []byte(password))

	salt := make([]byte, 0, len(c.Salt))
	salt = append(salt, c.Salt...)
	salt = append(salt, 0x00, 0x00, 0x00, 0x01)

	mac.Write(salt)
	result := mac.Sum(nil)

	prev := make([]byte, 0, len(result))
	prev = append(prev, result...)

	for i := 0; i < c.Iterations; i++ {
		mac.Write(prev)
		tmp := mac.Sum(nil)

		result = byteXOR(result, tmp)

		prev = tmp
	}

	return result
}

func GetServerSignatureMessage(sig []byte) string {
	return fmt.Sprintf("v=%s", base64.StdEncoding.EncodeToString(sig))
}

func GetClientKey(salted_pwd []byte) []byte {
	// Get Client Key
	mac := hmac.New(sha1.New, salted_pwd)
	mac.Write([]byte("Client Key"))
	return mac.Sum(nil)
}

func GetServerKey(salted_pwd []byte) []byte {
	// Get Server Key
	mac := hmac.New(sha1.New, salted_pwd)
	mac.Write([]byte("Server Key"))
	return mac.Sum(nil)
}

func GetServerSignature(auth string, serverk []byte) []byte {
	// Get Server Signature
	ssmac := hmac.New(sha1.New, serverk)
	ssmac.Write([]byte(auth))
	return ssmac.Sum(nil)
}

func GetClientSignature(auth string, storek []byte) []byte {
	// Get Client Signature
	skmac := hmac.New(sha1.New, storek)
	skmac.Write([]byte(auth))
	return skmac.Sum(nil)
}

func byteXOR(left, right []byte) []byte {
	res := make([]byte, len(left))
	for i := range left {
		res[i] = left[i] ^ right[i]
	}
	return res
}

func (r *Response) GetAuthMessage(c *Challenge, ca *ClientAuth) string {
	return fmt.Sprintf("n=%s,r=%s,%s,c=%s,r=%s", ca.UserName, base64.StdEncoding.EncodeToString(ca.Nonce),
		c.String(),
		r.Binding, base64.StdEncoding.EncodeToString(r.Nonce))
}

func (r *Response) GenerateProof(c *Challenge, ca *ClientAuth, password string) []byte {
	log.Println("Password", password)
	salted_pwd := c.SaltPassword(password)

	// Get Client and Server Keys
	clientk := GetClientKey(salted_pwd)
	log.Println("Client key", clientk)

	// Get Stored Key
	storek := sha1.Sum(clientk)

	// Build Auth Message
	auth := r.GetAuthMessage(c, ca)

	log.Println("Auth string", auth)

	client_sig := GetClientSignature(auth, storek[:])

	// Generate Proof
	client_proof := byteXOR(client_sig, clientk)

	r.server_sig = GetServerSignature(auth, GetServerKey(salted_pwd))

	return client_proof
}

func init() {
	auth.AddMechanism("SCRAM-SHA-1", func(e *auth.AuthElement, strm *stream.Stream) error {
		var state *SHAState
		if err := strm.State.Get(&state); err != nil {
			log.Println("SCRAM-SHA-1 is not available but tried to be used")
			return err
		}

		auth_data, err := mechanisms.DecodeBase64(e.Data, strm)
		if err != nil {
			return err
		}

		client_auth, err := NewClientAuthFromData(string(auth_data))
		if err != nil {
			return err
		}

		handler, err := newSHAHandler(strm, client_auth)
		if err != nil {
			return err
		}

		return handler.Handle()
	})

	auth.MechanismsElement.AddElement(mechanisms.NewMechanismElement(SHAElement("SCRAM-SHA-1")))
}
