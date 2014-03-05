package bind

/*
import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)
type bindElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	features.Elements
}

func (self *bindElement) IsRequiredFor(state interface{}) bool {
	return !state.(interface {
		Bound() bool
	}).Bound()
}

func (self *bindElement) CopyIfAvailable(state interface{}) interface{} {
	if self.IsRequiredFor(state) && state.(interface {
		Authenticated() bool
	}).Authenticated() {
		return self
	}
	return nil
}*/
