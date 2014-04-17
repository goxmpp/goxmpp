package mechanisms

import (
	"encoding/base64"
	"encoding/xml"
	"log"

	"github.com/dotdoom/goxmpp/extensions/features/auth"
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

type ChalengeElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl challenge"`
	Data    string   `xml:",chardata"`
}

func NewChalengeElement(data string) ChalengeElement {
	return ChalengeElement{Data: base64.StdEncoding.EncodeToString([]byte(data))}
}

type SuccessElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
}

type ResponseElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl response"`
	Data    string   `xml:",chardata"`
}

func NewResponseElement(data string) ResponseElement {
	return ResponseElement{Data: base64.StdEncoding.EncodeToString([]byte(data))}
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
	Method  Method   `xml:",chardata"`
}

type Method interface {
	IsAvailable(*stream.Stream) bool
}

func NewMechanismElement(method Method) *MechanismElement {
	return &MechanismElement{Method: method}
}

func (self *MechanismElement) CopyIfAvailable(strm *stream.Stream) elements.Element {
	if self.Method.IsAvailable(strm) {
		return self
	}
	return nil
}

func DecodeBase64(data string, strm *stream.Stream) ([]byte, error) {
	raw_data, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		log.Println("Could not decode Base64 in DigestMD5 handler:", err)
		if err := strm.WriteElement(auth.NewFailute(IncorrectEncoding{})); err != nil {
			return raw_data, err
		}
	}

	return raw_data, err
}
