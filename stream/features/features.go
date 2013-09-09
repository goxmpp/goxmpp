package features

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"

type Features struct {
	XMLName xml.Name `xml:"stream:features"`
	stream.InnerElements
}
