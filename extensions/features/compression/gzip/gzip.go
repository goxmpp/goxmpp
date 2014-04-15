package gzip

import (
	"compress/gzip"
	"io"
	"log"

	"github.com/dotdoom/goxmpp/extensions/features/compression"
	"github.com/dotdoom/goxmpp/stream"
	"github.com/dotdoom/goxmpp/stream/elements"
)

func init() {
	compression.CompressTemplate.AddElement(&compressor{BaseCompressor: compression.NewBaseCompressor("gzip")})
}

type State struct {
	Level int
}

type compressor struct {
	compression.BaseCompressor
}

func (c *compressor) CopyIfAvailable(s *stream.Stream) elements.Element {
	log.Println("Enabling compressor", c.Name())
	if c.IsAvailable(s) {
		return &compressor{BaseCompressor: compression.NewBaseCompressor(c.Name())}
	}
	return nil
}

func (c *compressor) GetReader(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

func (c *compressor) GetWriter(w io.Writer) io.WriteCloser {
	return gzip.NewWriter(w)
}
