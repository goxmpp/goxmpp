package compression

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/xep"

const (
	STREAM_NS = "urn:ietf:params:xml:ns:xmpp-sasl"
	STREAD_NODE = "compression"
)

func init(){
	stream.Registrator.Register(STREAM_NS + " " + STREAD_NODE, &CompressionHandler{
		Registrator: Registrator,
	})
}

var Registrator = xep.NewRegistrator()

type BaseCompression struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl compression"`
}

// This struct is used for marshaling
type Compression struct {
	BaseCompression
	stream.InnerElements
}

// This struct is used for unmarshaling and stream handling
type CompressionHandler {
	BaseCompression
	stream.InnerXML
}

func (self *CompressionHandler) Handle() {
// TODO: here will be handling and parsing code
}
