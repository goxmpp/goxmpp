package session

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"

type Feature struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-session session"`
	stream.InnerElements
}

func (self *Feature) CopyAvailableElements(sw *StreamWrapper) StreamFeature {
	if sw.State["authenticated"] == nil {
		return nil
	}
	return self.ExposeSubfeaturesTo(sw, new(SessionStreamFeature))
}
