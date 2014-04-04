package gzip

import (
	"compress/gzip"
	"io"

	"github.com/dotdoom/goxmpp/extensions/features/compression"
)

func init() {
	compression.AddMethod("gzip", &compressor{})
}

type State struct {
	Level int
}

type compressor struct{}

func (c *compressor) GetReader(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

func (c *compressor) GetWriter(w io.Writer) io.WriteCloser {
	return gzip.NewWriter(w)
}
