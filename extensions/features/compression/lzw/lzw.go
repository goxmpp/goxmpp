package lzw

import (
	"compress/lzw"
	"io"

	"github.com/dotdoom/goxmpp/extensions/features/compression"
)

func init() {
	compression.AddMethod("lzw", &compressor{})
}

type compressor struct{}

func (c *compressor) GetReader(r io.Reader) (io.ReadCloser, error) {
	return lzw.NewReader(r, lzw.LSB, 8), nil
}

func (c *compressor) GetWriter(w io.Writer) io.WriteCloser {
	return lzw.NewWriter(w, lzw.LSB, 8)
}
