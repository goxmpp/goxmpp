package auth

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type mechanismsElement struct {
	XMLName    xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Mechanisms []string `xml:"mechanism"`
	*features.Container
}

func newMechanismsElement() *mechanismsElement {
	return &mechanismsElement{
		Container: features.NewContainer(),
	}
}

func (self *mechanismsElement) IsRequiredFor(fs *features.State) bool {
	var state *State
	err := fs.Get(&state)
	return err != nil || state.UserName == ""
}

func (self *mechanismsElement) CopyIfAvailable(fs *features.State) elements.Element {
	if self.IsRequiredFor(fs) {
		x := newMechanismsElement()
		Features.CopyAvailableFeatures(fs, x.Container)
		return x
	}
	return nil
}

var Features = newMechanismsElement()
