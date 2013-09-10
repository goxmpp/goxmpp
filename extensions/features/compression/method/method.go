package method

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream/stanza"

type Method struct {
	XMLName xml.Name `xml:"method"`
	Name    string   `xml:",chardata"`
	stanza.InnerElements
}
