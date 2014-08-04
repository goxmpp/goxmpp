package stream

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"

	"github.com/goxmpp/xtream"
)

type StreamHandler func(ServerStream) error
type RawConfig map[string]json.RawMessage

var StreamXMLName = xml.Name{Local: "stream:stream"}

type Stream interface {
	xtream.Factory

	ID() string
	To() string
	From() string
	Open(StreamHandler) error
	Close() error
	State() *State
	Opened() bool
	UpdateRW(SwapRW) error
	ReadElement() (xtream.Element, error)
	WriteElement(xtream.Element) error
}

type ServerStream interface {
	Stream
	FeatureContainable
	SendFeatures() error
	ReOpen()
	Config() RawConfig
	SetDefaultNamespace(string)
	SetTo(string)
}

type serverStream struct {
	*stream
	*FeatureContainer
	reOpen bool
	config RawConfig
}

func NewServerStream(rw io.ReadWriteCloser, depGraph DependencyManageable, conf RawConfig) ServerStream {
	stream := NewStream(rw)
	stream.Factory = newStreamElementFactory()

	return &serverStream{
		stream:           stream,
		FeatureContainer: NewFeatureContainer(depGraph),
		config:           conf,
	}
}

func (s *serverStream) SendFeatures() error {
	return s.streamEncoder.Encode(s.FeatureContainer)
}

func (s *serverStream) ReOpen() {
	s.reOpen = true
}

func (s *serverStream) Config() RawConfig {
	return s.config
}

func (s *serverStream) Open(handler StreamHandler) error {
	if err := s.readSentOpen(); err != nil {
		return err
	}

	for s.HasRequired() {
		if s.reOpen {
			if err := s.readSentOpen(); err != nil {
				return err
			}
		}

		if err := handler(s); err != nil {
			return err
		}
	}

	s.stream.opened = true

	return nil
}

func (s *serverStream) readSentOpen() error {
	if err := s.stream.readOpen(); err != nil {
		return err
	}

	if !s.reOpen {
		if _, err := io.WriteString(s.rw, xml.Header); err != nil {
			return err
		}
	}

	if err := s.stream.sendOpen(); err != nil {
		return err
	}

	if err := s.SendFeatures(); err != nil {
		return err
	}

	s.reOpen = false
	return nil
}

func (s *serverStream) SetTo(to string) {
	s.stream.to = to
}

func (s *serverStream) SetDefaultNamespace(ns string) {
	s.stream.DefaultNamespace = ns
}

type stream struct {
	XMLName          xml.Name
	id               string `xml:"id,attr"`
	from             string `xml:"from,attr,omitempty"` // This holds server domain name.
	to               string `xml:"to,attr,omitempty"`   // This holds user JID after bind.
	version          string `xml:"version,attr"`
	DefaultNamespace string `xml:"-"`
	opened           bool
	state            *State
	xtream.Factory   `xml:"-"`
	Connection
}

func NewStream(rw io.ReadWriteCloser) *stream {
	st := &stream{Factory: xtream.NodeFactory, state: &State{}}
	st.SetRW(rw)
	return st
}

func (s *stream) ID() string {
	return s.id
}

func (s *stream) To() string {
	return s.to
}

func (s *stream) From() string {
	return s.from
}

func (s *stream) State() *State {
	return s.state
}

func (s *stream) Opened() bool {
	return s.opened
}

func (s *stream) sendClose() error {
	return s.streamEncoder.EncodeToken(xml.EndElement{xml.Name{Local: "stream:stream"}})
}

func (self *stream) readOpen() error {
	for {
		t, err := self.streamDecoder.Token()
		if err != nil {
			return err
		}
		switch t := t.(type) {
		case xml.ProcInst:
			// Good.
		case xml.StartElement:
			if t.Name.Local == "stream" {
				self.XMLName = t.Name
				self.XMLName.Local = "stream:stream"
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "to":
						self.from = attr.Value
					case "from":
						self.to = attr.Value
					case "version":
						self.version = attr.Value
					}
				}
				log.Printf("got <stream> from: %v, to: %v, version: %v\n", self.from, self.to, self.version)
				return nil
			}
		}
	}
}

func (s *stream) sendOpen() error {
	var start xml.StartElement
	start.Name = xml.Name{Local: "stream:stream", Space: "jabber:client"}
	start.Attr = append(start.Attr,
		xml.Attr{Name: xml.Name{Local: "xmlns:stream"}, Value: "http://etherx.jabber.org/streams"},
		xml.Attr{Name: xml.Name{Local: "xmlns:xml"}, Value: "http://www.w3.org/XML/1998/namespace"},
		xml.Attr{Name: xml.Name{Local: "to"}, Value: s.to},
		xml.Attr{Name: xml.Name{Local: "version"}, Value: s.version},
		xml.Attr{Name: xml.Name{Local: "from"}, Value: s.from},
	)
	if err := s.streamEncoder.EncodeToken(start); err != nil {
		return err
	}

	// xml.Encoder doesn't flush until it generated end tag
	// so we flush here to make it send stream's open tag
	return s.streamEncoder.Flush()
}

func (s *stream) Open() error {
	if err := s.sendOpen(); err != nil {
		return err
	}

	s.opened = true

	return s.readOpen()
}

func (self *stream) Close() error {
	if err := self.sendClose(); err != nil {
		return err
	}

	return self.Connection.Close()
}

func (self *stream) WriteElement(element xtream.Element) error {
	err := self.streamEncoder.Encode(element)
	if err != nil {
		log.Println("Error sending rely:", err)
	}
	return err
}

func (self *stream) ReadElement() (xtream.Element, error) {
	var err error
	var token xml.Token

	for token, err = self.streamDecoder.Token(); err == nil; token, err = self.streamDecoder.Token() {
		if start, ok := token.(xml.StartElement); ok {
			log.Printf("got element: %v (ns %v)\n", start.Name.Local, start.Name.Space)

			element := self.Factory.Get(&StreamXMLName, &start.Name)
			if element == nil {
				return nil, fmt.Errorf("Unknown node encountered: %s", start.Name.Local)
			}

			err := self.streamDecoder.DecodeElement(element, &start)
			return element, err
		}
	}

	return nil, err
}

type streamElementFactory struct {
	featuresFactory xtream.Factory
	elementsFactory xtream.Factory
}

func newStreamElementFactory() *streamElementFactory {
	return &streamElementFactory{xtream.NewFactory(), xtream.NodeFactory}
}

func (sef *streamElementFactory) Add(cons xtream.Constructor) {
	sef.featuresFactory.Add(cons)
}

func (sef *streamElementFactory) AddNamed(cons xtream.Constructor, outer, inner xml.Name) {
	sef.featuresFactory.AddNamed(cons, outer, inner)
}

func (sef *streamElementFactory) Get(outer, inner *xml.Name) xtream.Element {
	setFactory := func(el xtream.Element) xtream.Element {
		if innerEl, ok := el.(xtream.Registrable); ok {
			innerEl.SetFactory(sef)
		}
		return el
	}

	if e := sef.featuresFactory.Get(outer, inner); e != nil {
		return setFactory(e)
	}

	return setFactory(sef.elementsFactory.Get(outer, inner))
}
