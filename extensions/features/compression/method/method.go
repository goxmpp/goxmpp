package method

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/extensions/features/compression"
	"github.com/dotdoom/goxmpp/stream"
)

func init() {
	compression.HandlerRegistrator.Register(" method", func() stream.Element { return &BasicMethod{} })
}

type BasicMethod struct {
	XMLName xml.Name `xml:"method"`
	Name    string   `xml:",chardata"`
}

type Method struct {
	BasicMethod
	stream.InnerElements
}
