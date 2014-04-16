package mechanisms

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
)

func init() {
	stream.StreamFactory.AddConstructor(func() elements.Element {
		return &ResponseElement{}
	})
	stream.StreamFactory.AddConstructor(func() elements.Element {
		return &Abort{}
	})
}

type SuccessElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
}

type ResponseElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl response"`
	Data    string   `xml:",chardata"`
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
