package lzw

import (
	"compress/lzw"
	"io"
	"log"

	"github.com/goxmpp/goxmpp/extensions/features/compression"
	"github.com/goxmpp/goxmpp/stream"
	"github.com/goxmpp/goxmpp/stream/elements"
)

func init() {
	compression.CompressTemplate.AddElement(&compressor{BaseCompressor: compression.NewBaseCompressor("lzw")})
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
	return lzw.NewReader(r, lzw.LSB, 8), nil
}

func (c *compressor) GetWriter(w io.Writer) io.WriteCloser {
	return lzw.NewWriter(w, lzw.LSB, 8)
}
