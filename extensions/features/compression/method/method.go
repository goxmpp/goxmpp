package method

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/extensions/features/compression"
	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/stanza"
)

func init() {
	compression.HandlerRegistrator.Register(" method", &BasicMethod{})
}

type BasicMethod struct {
	XMLName xml.Name `xml:"method"`
	Name    string   `xml:",chardata"`
}

type Method struct {
	BasicMethod
	stanza.InnerElements
}

func (self *BasicMethod) HandleElement(sw *stream.Wrapper) {

}
