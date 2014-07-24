package auth

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"log"

	"github.com/goxmpp/goxmpp/extensions/features/auth/mechanisms"
	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/features"
)

type MechanismsElement struct {
	XMLName    xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Mechanisms []*mechanisms.MechanismElement
}

func newMechanismsElement(opts features.Options) features.BasicFeature {
	auth := &MechanismsElement{
		Mechanisms: make([]*mechanisms.MechanismElement, 0),
	}

	for mech := range mechanism_handlers {
		auth.Mechanisms = append(auth.Mechanisms, mechanisms.NewMechanismElement(mech))
	}

	return auth
}

func (me *MechanismsElement) NewHandler() features.FeatureHandler {
	return &AuthElement{}
}

var mechanism_handlers map[string]Handler = make(map[string]Handler)

func AddMechanism(name string, handler Handler) {
	mechanism_handlers[name] = handler
}

type Handler func(*AuthElement, *stream.Stream) error

type AuthElement struct {
	XMLName   xml.Name `xml:"auth"`
	Mechanism string   `xml:"mechanism,attr"`
	Data      string   `xml:",chardata"`
}

func (self *AuthElement) Handle(strm features.FeatureContainable, opts features.Options) error {
	st := strm.(*stream.Stream)

	if handler := mechanism_handlers[self.Mechanism]; handler != nil {
		if err := handler(self, st); err != nil {
			log.Println("Authorization failed:", err)
			if err := st.WriteElement(NewFailute(NotAuthorized{})); err != nil {
				return err
			}
			return err
		}
	} else {
		if err := st.WriteElement(NewFailute(InvalidMechanism{})); err != nil {
			return err
		}
		return fmt.Errorf("No handler for mechanism %v", self.Mechanism)
	}

	return nil
}

func DecodeBase64(data string, strm *stream.Stream) ([]byte, error) {
	raw_data, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		log.Println("Could not decode Base64 in DigestMD5 handler:", err)
		if err := strm.WriteElement(NewFailute(mechanisms.IncorrectEncoding{})); err != nil {
			return raw_data, err
		}
	}

	return raw_data, err
}

func init() {
	features.FeatureFactory.Add("auth", &features.FeatureFactoryElement{
		Constructor: func(opts features.Options) *features.Feature {
			return features.NewFeature("auth", newMechanismsElement(opts), true)
		},
		Name:   xml.Name{Local: "auth"},
		Parent: stream.StreamXMLName,
	})
}
