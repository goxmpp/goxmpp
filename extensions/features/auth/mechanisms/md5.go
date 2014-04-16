package mechanisms

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/dotdoom/goxmpp/extensions/features/auth"
	"github.com/dotdoom/goxmpp/stream"
)

const A2_AUTH_SUFFIX = "00000000000000000000000000000000"

type ChalengeElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl challenge"`
	Data    string   `xml:",chardata"`
}

func NewChalengeElement(data string) ChalengeElement {
	return ChalengeElement{Data: base64.StdEncoding.EncodeToString([]byte(data))}
}

type Chalenge struct {
	Realm     []string
	Nonce     string
	QOP       string
	Charset   string
	Algorithm string
}

func NewChalenge(realm []string) (*Chalenge, error) {
	nonce := make([]byte, 14)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return &Chalenge{
		Realm:     realm,
		Nonce:     fmt.Sprintf("%x", nonce),
		QOP:       "auth",
		Algorithm: "md5-sess",
	}, nil
}

func (c *Chalenge) String() string {
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

func decodeMD5Response(data []byte) *Response {
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

type DigestMD5State struct {
	ValidateMD5 func(*Chalenge, *Response) bool
	Realm       []string
	Host        string
}

func init() {
	auth.AddMechanism("DIGEST-MD5", func(e *auth.AuthElement, strm *stream.Stream) error {
		var md5_state *DigestMD5State
		if err := strm.State.Get(&md5_state); err != nil {
			return err
		}
		// TODO Need to handle aborts

		// First send chalenge with nonce
		//chalenge := `realm="cataclysm.cx",nonce="OA6MG9tEQGm2hh",qop="auth",charset=utf-8,algorithm=md5-sess`
		chalenge, err := NewChalenge(md5_state.Realm)
		if err != nil {
			log.Println("Could not create a chalenge")
			return err
		}
		if err := strm.WriteElement(NewChalengeElement(chalenge.String())); err != nil {
			return err
		}

		// Receive a response with encoded MD5
		el, err := strm.ReadElement()
		if err != nil {
			return err
		}

		resp_el, ok := el.(*ResponseElement)
		if !ok || resp_el.Data == "" {
			return errors.New("Wrong response received")
		}

		// Check MD5
		raw_resp_data, err := base64.StdEncoding.DecodeString(resp_el.Data)

		log.Println("Sent chalenge", chalenge.String())
		log.Println("Received response", string(raw_resp_data))
		if err != nil {
			log.Println("Could not decode Base64 in DigestMD5 handler:", err)
			if err := strm.WriteElement(auth.NewFailute(IncorrectEncoding{})); err != nil {
				return err
			}
			return err
		}

		resp := decodeMD5Response(raw_resp_data)
		log.Printf("Chalenge object %#v", chalenge)
		log.Printf("Response object %#v", resp)

		if err := basicResponseCheck(chalenge, resp, md5_state); err != nil {
			return err
		}
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
			// Need to send meaningful error to other side
			return errors.New("Wrong response received")
		}

		if err := strm.WriteElement(SuccessElement{}); err != nil {
			return err
		}

		var auth_state *auth.AuthState
		if err := strm.State.Get(&auth_state); err != nil {
			auth_state = &auth.AuthState{}
			strm.State.Push(auth_state)
		}

		auth_state.UserName = resp.UserName

		strm.ReOpen = true

		return nil
	})

	auth.MechanismsElement.AddElement(newMechanismElement("DIGEST-MD5"))
}

func GenerateResponseHash(c *Chalenge, r *Response, password string) string {
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

func basicResponseCheck(c *Chalenge, r *Response, state *DigestMD5State) error {
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
