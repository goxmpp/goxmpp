package method

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/extensions/features/compression"
	"github.com/dotdoom/goxmpp/stream/elements"
)

func init() {
	compression.CompressionFactory.AddConstructor(" method", func() elements.Element { return &BasicMethod{} })
}

type BasicMethod struct {
	XMLName xml.Name `xml:"method"`
	Name    string   `xml:",chardata"`
}

type Method struct {
	BasicMethod
	elements.InnerElements
}
