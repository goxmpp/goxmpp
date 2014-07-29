package starttls

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net"

	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/features"
)

func init() {
	features.FeatureFactory.Add(
		"starttls",
		&features.FeatureFactoryElement{
			Constructor: func(opts features.Options) *features.Feature {
				return features.NewFeature("starttls", NewStartTLSFeature(false), false)
			},
			Name:   xml.Name{Local: "starttls", Space: "urn:ietf:params:xml:ns:xmpp-tls"},
			Parent: stream.StreamXMLName,
		},
	)
}

type StartTLSFeatureElement struct {
	Required        bool `xml:"required,omitempty"`
	StartTLSElement      // will get XMLName from here
}

func NewStartTLSFeature(required bool) *StartTLSFeatureElement {
	return &StartTLSFeatureElement{Required: required}
}

func (tls *StartTLSFeatureElement) NewHandler() features.FeatureHandler {
	return &StartTLSElement{}
}

type StartTLSElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
}

type TLSConfig struct {
	PEMPath string
	KeyPath string
}

func NewTLSConfig(pem, key string) *TLSConfig {
	return &TLSConfig{PEMPath: pem, KeyPath: key}
}

type StartTLSState struct {
	Required bool
	Started  bool
	Config   *TLSConfig
}

func NewStartTLSState(required bool, conf *TLSConfig) *StartTLSState {
	return &StartTLSState{
		Started:  false,
		Config:   conf,
		Required: required,
	}
}

func (s *StartTLSElement) Handle(st *stream.Stream, opts features.Options) error {
	var state *StartTLSState
	if err := st.State.Get(&state); err != nil {
		return err
	}

	cert, err := tls.LoadX509KeyPair(state.Config.PEMPath, state.Config.KeyPath)
	if err != nil {
		log.Println("Could not load keys:", err)
		return err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.VerifyClientCertIfGiven,
		Rand:         rand.Reader,
	}

	err = st.UpdateRW(func(srwc io.ReadWriteCloser) (io.ReadWriteCloser, error) {
		if conn, ok := srwc.(net.Conn); ok {
			tls_conn := tls.Server(conn, config)

			// Once we inialized - let client proceed
			if err := st.WriteElement(&ProceedElement{}); err != nil {
				return nil, err
			}

			// Now do a handshake
			if err := tls_conn.Handshake(); err != nil {
				log.Println("TLS Handshake error:", err)
				return nil, err
			}
			return tls_conn, nil
		}
		return nil, errors.New("Wrong ReadWriteCloser, expected connection")
	})
	if err != nil {
		log.Println("Could not replace connection", err)
		return err
	}

	state.Started = true
	st.ReOpen = true

	return nil
}

type ProceedElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}
