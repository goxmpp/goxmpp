package compression

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"

const (
	STREAM_NS   = "urn:ietf:params:xml:ns:xmpp-sasl"
	STREAD_NODE = "compression"
)

func init() {
	stream.HandlerRegistrator.Register(STREAM_NS+" "+STREAD_NODE, func() stream.Element {
		return &CompressionHandler{InnerXML: stream.InnerXML{Registrator: HandlerRegistrator}}
	})
}

var HandlerRegistrator = stream.NewElementGeneratorRegistrator()

type BaseCompression struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl compression"`
}

// This struct is used for marshaling
type CompressionFeature struct {
	BaseCompression
	stream.InnerElements
}

// This struct is used for unmarshaling and stream handling
type CompressionHandler struct {
	BaseCompression
	stream.InnerXML
}
