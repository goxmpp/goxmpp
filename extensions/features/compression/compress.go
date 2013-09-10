package compression

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/decoder"
import "github.com/dotdoom/goxmpp/stream/stanza"
import "github.com/dotdoom/goxmpp/extensions"

const (
	STREAM_NS = "urn:ietf:params:xml:ns:xmpp-sasl"
	STREAD_NODE = "compression"
)

func init(){
	stream.Registrator.Register(STREAM_NS + " " + STREAD_NODE, &CompressionHandler{
		Registrator: HandlerRegistrator,
	})
}

var HandlerRegistrator = extensions.NewHandlerRegistrator()

type BaseCompression struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl compression"`
}

// This struct is used for marshaling
type CompressionFeature struct {
	BaseCompression
	stanza.InnerElements
}

// This struct is used for unmarshaling and stream handling
type CompressionHandler {
	BaseCompression
	decoder.InnerXML
}

func (self *CompressionHandler) HandleElement() {
// TODO: here will be handling and parsing code
}
