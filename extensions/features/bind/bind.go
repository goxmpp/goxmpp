package bind

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

type bindElement struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	features.Elements
}

func (self *bindElement) IsRequiredFor(fs features.State) bool {
	return fs["bound"] == nil
}

func (self *bindElement) CopyIfAvailable(fs features.State) interface{} {
	if self.IsRequiredFor(fs) && fs["authenticated"] != nil {
		return self.CopyAvailableFeatures(fs, new(bindElement))
	}
	return nil
}

func init() {
	features.List.AddElement(new(bindElement))
}
