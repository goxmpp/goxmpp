package compression

import (
	"encoding/xml"
	"github.com/dotdoom/goxmpp/stream/elements"
)

const (
	STREAM_NS   = "urn:ietf:params:xml:ns:xmpp-sasl"
	STREAM_NODE = "compression"
)

func init() {
	elements.GlobalFeaturesFactory.AddConstructor(STREAM_NS+" "+STREAM_NODE, func() elements.Element {
		return &CompressionHandler{
			UnmarshallableElements: elements.UnmarshallableElements{ElementFactory: ElementFactory}}
	})
}

var ElementFactory = elements.NewFactory()

type BaseCompression struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl compression"`
}

// This struct is used for marshaling
type CompressionFeature struct {
	BaseCompression
	elements.Elements
}

// This struct is used for unmarshaling and stream handling
type CompressionHandler struct {
	BaseCompression
	elements.UnmarshallableElements
}
