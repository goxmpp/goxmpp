package zlib

import (
	"compress/zlib"
	"io"

	"github.com/goxmpp/goxmpp/extensions/features/compression"
)

func init() {
	compression.Methods["zlib"] = compressor{compression.BaseCompressor{MethodName: "zlib"}}
}

type compressor struct {
	compression.BaseCompressor
}

func (c compressor) GetReader(r io.Reader) (io.ReadCloser, error) {
	return zlib.NewReader(r)
}

func (c compressor) GetWriter(w io.Writer) io.WriteCloser {
	return zlib.NewWriter(w)
}
