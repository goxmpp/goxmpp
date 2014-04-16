package compression

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"

	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
	"github.com/dotdoom/goxmpp/stream/elements/features"
)

func init() {
	features.Tree.AddElement(CompressTemplate)
	stream.StreamFactory.AddConstructor(func() elements.Element {
		return NewCompressHandler()
	})
}

type Compressor interface {
	GetReader(io.Reader) (io.ReadCloser, error)
	GetWriter(io.Writer) io.WriteCloser
	Name() string
	IsAvailable(*stream.Stream) bool
}
type CompressorConfig struct {
	Level int
}
type CompressState struct {
	Compressed bool
	Config     map[string]CompressorConfig
}

func NewCompressState() *CompressState {
	return &CompressState{Compressed: false}
}

type BaseCompressor struct {
	XMLName        xml.Name `xml:"method"`
	CompressorName string   `xml:",chardata"`
}

func NewBaseCompressor(name string) BaseCompressor {
	return BaseCompressor{CompressorName: name}
}

func (bc BaseCompressor) Name() string {
	return bc.CompressorName
}

func (bc BaseCompressor) IsAvailable(stream *stream.Stream) bool {
	// TODO Add some logic to check if this method is available
	return true
}

// This struct is used for marshaling
type compression struct {
	XMLName xml.Name `xml:"http://jabber.org/features/compress compression"`
	*features.Container
}

func NewCompression() *compression {
	return &compression{
		Container: features.NewContainer(),
	}
}

var CompressTemplate = NewCompression()

func (c *compression) CopyIfAvailable(stream *stream.Stream) elements.Element {
	var state *CompressState
	err := stream.State.Get(&state)

	if err != nil || state.Compressed {
		return nil
	}

	compress := NewCompression()
	c.CopyAvailableFeatures(stream, compress.Container)
	return compress
}

type compressElement struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/compress compress"`
	Method  string   `xml:"method"`
}

func NewCompressHandler() *compressElement {
	return &compressElement{}
}

func (c *compressElement) Handle(s *stream.Stream) error {
	var compressor Compressor

	for _, element := range CompressTemplate.Elements() {
		if compr, ok := element.(Compressor); ok && compr.Name() == c.Method && compr.IsAvailable(s) {
			compressor = compr
			break
		}
	}

	if compressor == nil {
		if err := s.WriteElement(&MethodNotSupportedError{}); err != nil {
			return err
		}
		return fmt.Errorf("Unsupported compression method requested")
	}

	var state *CompressState
	if err := s.State.Get(&state); err != nil {
		if err := s.WriteElement(&ProcessingFailedError{}); err != nil {
			return err
		}
		return err
	}

	state.Compressed = true
	if err := s.WriteElement(&CompressionSuccess{}); err != nil {
		return err
	}

	if err := swapStreamRW(s, compressor); err != nil {
		if err := s.WriteElement(&ProcessingFailedError{}); err != nil {
			return err
		}
		return err
	}

	s.ReOpen = true
	return nil
}

func swapStreamRW(strm *stream.Stream, compressor Compressor) error {
	return strm.UpdateRW(
		func(srwc io.ReadWriteCloser) (io.ReadWriteCloser, error) {
			writer := compressor.GetWriter(srwc)
			reader, err := compressor.GetReader(srwc)
			if err != nil {
				log.Println("Could not create compressed reader", err)
				if err := strm.WriteElement(&SetupFailedError{}); err != nil {
					return nil, err
				}
				return nil, err
			}

			return NewCompressionReadWriter(srwc, reader, writer), nil
		})
}

type CompressionSuccess struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/compress compressed"`
}

type MethodNotSupportedError struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/compress failure"`
	Error   xml.Name `xml:"unsupported-method"`
}

type SetupFailedError struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/compress failure"`
	Error   xml.Name `xml:"setup-failed"`
}

type ProcessingFailedError struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/compress failure"`
	Error   xml.Name `xml:"processing-failed"`
}

type compressedReadWriter struct {
	source io.ReadWriteCloser
	reader io.ReadCloser
	writer io.WriteCloser
}
type Flusher interface {
	Flush() error
}

func NewCompressionReadWriter(s io.ReadWriteCloser, r io.ReadCloser, w io.WriteCloser) *compressedReadWriter {
	return &compressedReadWriter{source: s, reader: r, writer: w}
}

func (c *compressedReadWriter) Read(b []byte) (int, error) {
	return c.reader.Read(b)
}

func (c *compressedReadWriter) Write(b []byte) (int, error) {
	n, err := c.writer.Write(b)
	if err != nil {
		return n, err
	}
	// Need to flush here, otherwise data won't get to network
	// Data will be buffered on lower lever e.g. TCP
	if f, ok := c.writer.(Flusher); ok {
		return n, f.Flush()
	}

	return n, nil
}

func (c *compressedReadWriter) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	if err := c.writer.Close(); err != nil {
		return err
	}
	return c.source.Close()
}
