package mechanisms

import (
	"encoding/base64"
	"encoding/xml"
	"errors"

	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/xtream"
)

func init() {
	xtream.NodeFactory.Add(func(*xml.Name) xtream.Element {
		return &ResponseElement{}
	})
	xtream.NodeFactory.Add(func(*xml.Name) xtream.Element {
		return &Abort{}
	})
}

type ChallengeElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl challenge"`
	Data    string   `xml:",chardata"`
}

func NewChallengeElement(data []byte) ChallengeElement {
	return ChallengeElement{Data: base64.StdEncoding.EncodeToString(data)}
}

type SuccessElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
	Data    string   `xml:",chardata"`
}

func NewSuccessElement(data []byte) SuccessElement {
	return SuccessElement{Data: base64.StdEncoding.EncodeToString(data)}
}

type ResponseElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl response" parent:"stream:stream"`
	Data    string   `xml:",chardata"`
}

func NewResponseElement(data string) ResponseElement {
	return ResponseElement{Data: base64.StdEncoding.EncodeToString([]byte(data))}
}

func ReadResponse(strm *stream.Stream) (*ResponseElement, error) {
	el, err := strm.ReadElement()
	if err != nil {
		return nil, err
	}

	resp, ok := el.(*ResponseElement)
	if !ok {
		// Need to send meaningful error to other side
		return nil, errors.New("Wrong response received")
	}

	return resp, nil
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
	XMLName xml.Name `xml:"abort" parent:"stream:stream"`
}

type MechanismElement struct {
	XMLName   xml.Name `xml:"mechanism"`
	Mechanism string   `xml:",chardata"`
}

func NewMechanismElement(mech string) *MechanismElement {
	return &MechanismElement{Mechanism: mech}
}
