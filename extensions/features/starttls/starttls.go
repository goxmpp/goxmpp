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
				conf := opts.(*startTLSConf)
				return features.NewFeature("starttls", NewStartTLSFeature(conf), false, conf)
			},
			Name:   xml.Name{Local: "starttls", Space: "urn:ietf:params:xml:ns:xmpp-tls"},
			Parent: stream.StreamXMLName,
			Config: func() interface{} { return &startTLSConf{} },
		},
	)
}

type startTLSConf struct {
	Required bool   `json:"required"`
	PEMPath  string `json:"pem"`
	KeyPath  string `json:"key"`
}

type StartTLSFeatureElement struct {
	Required        bool `xml:"required,omitempty"`
	StartTLSElement      // will get XMLName from here
}

func NewStartTLSFeature(conf *startTLSConf) *StartTLSFeatureElement {
	return &StartTLSFeatureElement{Required: conf.Required}
}

func (tls *StartTLSFeatureElement) NewHandler() features.FeatureHandler {
	return &StartTLSElement{}
}

type StartTLSElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
}

type StartTLSState struct {
	Started bool
}

func NewStartTLSState() *StartTLSState {
	return &StartTLSState{
		Started: false,
	}
}

func (s *StartTLSElement) Handle(st stream.ServerStream, opts features.Options) error {
	conf := opts.(*startTLSConf)
	cert, err := tls.LoadX509KeyPair(conf.PEMPath, conf.KeyPath)
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

	state := NewStartTLSState()
	state.Started = true
	st.State().Push(state)
	st.ReOpen()

	return nil
}

type ProceedElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}
