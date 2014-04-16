package mechanisms

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/dotdoom/goxmpp/extensions/features/auth"
	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
)

type SuccessElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
}

type ResponseElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl response"`
	Data    string   `xml:",chardata,omitempty"`
}

type NotAuthorized struct {
	XMLName xml.Name `xml:"not-authorized"`
}

type IncorrectEncoding struct {
	XMLName xml.Name `xml:"incorrect-encoding"`
}

type InvalidAuthID struct {
	XMLName xml.Name `xml:"invalid-authzid"`
}

type MachanismTooWeak struct {
	XMLName xml.Name `xml:"mechanism-too-weak"`
}

type Aborted struct {
	XMLName xml.Name `xml:"aborted"`
}

type Abort struct {
	XMLName xml.Name `xml:"abort"`
}

type ChalengeElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl challenge"`
	Data    string   `xml:",chardata,omitempty"`
}

func NewChalengeElement(data string) ChalengeElement {
	return ChalengeElement{Data: base64.StdEncoding.EncodeToString([]byte(data))}
}

type PlainState struct {
	VerifyUserAndPassword func(string, string) bool
	RequireEncryption     bool
}

type Chalenge struct {
	Realm     []string
	Nonce     string
	QOP       string
	Charset   string
	Algorithm string
}

func (c *Chalenge) String() string {
	str := []string{fmt.Sprintf("nonce=%s", c.Nonce), fmt.Sprintf("algorithm=%s", c.Algorithm)}
	for _, realm := range c.Realm {
		str = append(str, fmt.Sprintf("realm=%s", realm))
	}
	if c.QOP != "" {
		str = append(str, fmt.Sprintf("qop=%s", c.QOP))
	}
	if c.Charset != "" {
		str = append(str, fmt.Sprintf("charset=%s", c.QOP))
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

func decodeMD5Response(data []byte) *Response {
	resp := &Response{}

	for _, param := range bytes.Split(data, []byte(",")) {
		key_val := bytes.SplitN(param, []byte("="), 2)
		val := string(key_val[1])
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

type DigestMD5State struct {
	ValidateMD5 func(*Chalenge, *Response) bool
	GetChalenge func() *Chalenge
}

type MechanismElement struct {
	XMLName xml.Name `xml:"mechanism"`
	Name    string   `xml:",chardata"`
}

func newMechanismElement(name string) *MechanismElement {
	return &MechanismElement{Name: name}
}

func (self *MechanismElement) CopyIfAvailable(strm *stream.Stream) elements.Element {
	switch self.Name {
	case "PLAIN":
		var plain_state *PlainState
		if err := strm.State.Get(&plain_state); err == nil {
			return self
		}
	case "DIGEST-MD5":
		var md5_state *DigestMD5State
		if err := strm.State.Get(&md5_state); err == nil {
			return self
		}
	}
	return nil
}

var usernamePasswordSeparator = []byte{0}

func init() {
	auth.AddMechanism("PLAIN", func(e *auth.AuthElement, stream *stream.Stream) error {
		b, err := base64.StdEncoding.DecodeString(e.Data)
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
	auth.AddMechanism("DIGEST-MD5", func(e *auth.AuthElement, strm *stream.Stream) error {
		var md5_state *DigestMD5State
		if err := stream.State.Get(&md5_state); err != nil {
			return err
		}
		// TODO Need to handle aborts

		// First send chalenge with nonce
		//chalenge := `realm="cataclysm.cx",nonce="OA6MG9tEQGm2hh",qop="auth",charset=utf-8,algorithm=md5-sess`
		chalenge := md5_state.GetChalenge()
		if err := strm.WriteElement(NewChalengeElement(chalenge.String())); err != nil {
			return err
		}

		// Receive a response with encoded MD5
		el, err := strm.ReadElement()
		if err != nil {
			return err
		}

		resp_el, ok := el.(*ResponseElement)
		if !ok || resp.Data == "" {
			return errors.New("Wrong response received")
		}

		// Check MD5
		raw_resp_data, err := _base64.StdEncoding.DecodeString(resp_el.Data)
		if err != nil {
			return err
		}
		resp := decodeMD5Response(raw_resp_data)
		if !md5_state.ValidateMD5(chalenge, resp) {

			return errors.New("AUTH FAILED")
		}

		// Send response
		if err := strm.WriteElement(NewChalengeElement("rspauth")); err != nil {
			return err
		}

		el, err = strm.ReadElement()
		if err != nil {
			return err
		}
		if resp, ok := el.(*ResponseElement); !ok || resp.Data != "" {
			// Need to send failure
			return errors.New("Wrong response received")
		}

		if err := strm.WriteElement(SuccessElement{}); err != nil {
			return err
		}

		auth_state.UserName = resp.username

		strm.ReOpen = true

		return nil
	})

	auth.MechanismsElement.AddElement(newMechanismElement("PLAIN"))
	auth.MechanismsElement.AddElement(newMechanismElement("DIGEST-MD5"))

	stream.StreamFactory.AddConstructor(func() elements.Element {
		return &ResponseElement{}
	})
	stream.StreamFactory.AddConstructor(func() elements.Element {
		return &Abort{}
	})
}
