package compression

import "encoding/xml"
import "github.com/dotdoom/goxmpp/stream"
import "github.com/dotdoom/goxmpp/stream/elements"

const (
	STREAM_NS   = "urn:ietf:params:xml:ns:xmpp-sasl"
	STREAD_NODE = "compression"
)

func init() {
	stream.GlobalFeaturesFactory.AddConstructor(STREAM_NS+" "+STREAD_NODE, func() elements.Element {
		return &CompressionHandler{InnerXML: stream.InnerXML{ElementFactory: ElementFactory}}
	})
}

var ElementFactory = elements.NewFactory()

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
