package gzip

import (
	"compress/gzip"
	"io"

	"github.com/goxmpp/goxmpp/extensions/features/compression"
)

func init() {
	compression.Methods["gzip"] = compressor{compression.BaseCompressor{MethodName: "gzip"}}
}

type State struct {
	Level int
}

type compressor struct {
	compression.BaseCompressor
}

func (c compressor) GetReader(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

func (c compressor) GetWriter(w io.Writer) io.WriteCloser {
	return gzip.NewWriter(w)
}
