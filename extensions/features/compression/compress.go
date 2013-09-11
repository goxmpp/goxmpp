package compression

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/stanza"

const (
	STREAM_NS   = "urn:ietf:params:xml:ns:xmpp-sasl"
	STREAD_NODE = "compression"
)

func init() {
	stream.HandlerRegistrator.Register(STREAM_NS+" "+STREAD_NODE, &CompressionHandler{
		InnerXML: stream.InnerXML{Registrator: HandlerRegistrator},
	})
}

var HandlerRegistrator = stream.NewElementHandlerRegistrator()

type BaseCompression struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl compression"`
}

// This struct is used for marshaling
type CompressionFeature struct {
	BaseCompression
	stanza.InnerElements
}

// This struct is used for unmarshaling and stream handling
type CompressionHandler struct {
	BaseCompression
	stream.InnerXML
}

func (self *CompressionHandler) HandleElement(sw *stream.Wrapper) {
	// TODO: here will be handling and parsing code
}
