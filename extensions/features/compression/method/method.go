package method

import "encoding/xml"
import "github.com/dotdoom/goxmpp/extensions"
import "github.com/dotdoom/goxmpp/stream/stanza"
import "github.com/dotdoom/goxmpp/stream/decoder"
import "github.com/dotdoom/goxmpp/extensions/features/compression"

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
