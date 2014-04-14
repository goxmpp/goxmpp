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
	features.Tree.AddElement(NewCompression())
	stream.StreamFactory.AddConstructor(func() elements.Element {
		return NewCompressionHandler()
	})
}

type CompressorConstructor func(*stream.Stream) (Compressor, error)
type Compressor interface {
	GetReader(io.Reader) (io.ReadCloser, error)
	GetWriter(io.Writer) io.WriteCloser
}

var compressionMethods map[string]CompressorConstructor = make(map[string]CompressorConstructor)

func AddMethod(name string, handler CompressorConstructor) {
	compressionMethods[name] = handler
}

func NewCompression() *compression {
	comp := &compression{
		Methods:   make([]string, 0, len(compressionMethods)),
		Container: features.NewContainer(),
	}
	for method := range compressionMethods {
		comp.Methods = append(comp.Methods, method)
	}
	return comp
}

// This struct is used for marshaling
type compression struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/compress compression"`
	Methods []string `xml:"method"`
	*features.Container
}

type CompressElement struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/compress compress"`
	Method  string   `xm;"method"`
}

func (c *CompressElement) Handle(stream *stream.Stream) error {
	if compressor, ok := compressionMethods[c.Method]; ok {
		stream.WriteElement(&CompressionSuccess{})
		if err := swapStreamRW(stream, compressor); err != nil {
			stream.WriteElement(&ProcessingFailedError{})
			return err
		}
		return nil
	}

	stream.WriteElement(&MethodNotSupportedError{})
	return fmt.Errorf("Unsupported compression method requested")
}

func swapStreamRW(strm *stream.Stream, compressor Compressor) error {
	var err error
	strm.UpdateRW(
		func(srwc io.ReadWriteCloser) io.ReadWriteCloser {
			var reader io.ReadCloser

			writer := compressor.GetWriter(srwc)
			reader, err = compressor.GetReader(srwc)
			if err != nil {
				log.Println(err)
				strm.WriteElement(&SetupFailedError{})
				return nil
			}

			return NewCompressionReadWriter(srwc, reader, writer)
		},
	)

	if err != nil {
		return err
	}

	return strm.ReadOpen()
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

func NewCompressionReadWriter(s io.ReadWriteCloser, r io.ReadCloser, w io.WriteCloser) *compressedReadWriter {
	return &compressedReadWriter{source: s, reader: r, writer: w}
}

func (c *compressedReadWriter) Read(b []byte) (int, error) {
	return c.reader.Read(b)
}

func (c *compressedReadWriter) Write(b []byte) (int, error) {
	return c.writer.Write(b)
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
