package compression

import (
	"encoding/xml"

	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

const (
	STREAM_NS   = "urn:ietf:params:xml:ns:xmpp-sasl"
	STREAM_NODE = "compression"
)

func init() {
	features.Factory.AddConstructor(STREAM_NS+" "+STREAM_NODE, func() elements.Element {
		return NewCompressionHandler()
	})
}

var ElementFactory = elements.NewElementFactory()

func NewCompressionHandler() *CompressionHandler {
	return &CompressionHandler{InnerElements: elements.NewInnerElements(ElementFactory)}
}

type BaseCompression struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl compression"`
}

// This struct is used for marshaling
type CompressionFeature struct {
	BaseCompression
	*elements.InnerElements
}

// This struct is used for unmarshaling and stream handling
type CompressionHandler struct {
	BaseCompression
	*elements.InnerElements
}
