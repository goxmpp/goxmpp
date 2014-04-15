package starttls

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net"

	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

func init() {
	features.Tree.AddElement(NewStartTLSFeature())
	stream.StreamFactory.AddConstructor(func() elements.Element {
		return &StartTLSElement{}
	})
}

type StartTLSFeatureElement struct {
	Required        bool `xml:"required,omitempty"`
	StartTLSElement      // will get XMLName from here
}

func NewStartTLSFeature(required bool) *StartTLSFeatureElement {
	return &StartTLSFeatureElement{Required: required}
}

func (s *StartTLSFeatureElement) CopyIfAvailable(s *stream.Stream) elements.Element {
	// Handcoded for now because this will be reworked anyway
	return NewStartTLSFeature(true)
}

type StartTLSElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
}

type TSLConfig struct {
	PEMPath string
	KeyPath string
}

type StartTLSState struct {
	Started bool
	Config  TLSConfig
}

func (s *StartTLSElement) Handle(s *stream.Stream) error {
	var state *StartTLSState
	if err := s.State.Get(&state); err != nil {
		return err
	}

	cert, err := tls.LoadX509KeyPair(state.Config.PEMPath, state.Config.KeyPath)
	if err != nil {
		log.Println("Could not load keys:", err)
		return err
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}

	err = s.UpdateRW(func(srwc io.ReadWriteCloser) (io.ReadWriteCloser, error) {
		if conn, ok := srwc.(net.Conn); ok {
			tls_conn, err := tls.Server(conn, config)
			if err != nil {
				return nil, err
			}
			s.WriteElement(&ProceedElement{})
			return tls_conn, nil
		}
		return nil, errors.New("Wrong ReadWriteCloser, expected connection")
	})
	if err != nil {
		log.Println("Could not replace connection", err)
	}

	return err
}

type ProceedElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}
