package method

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"

type Method struct {
	XMLName xml.Name `xml:"method"`
	Name    string   `xml:",chardata"`
	stream.InnerElements
}
