package lzw

import (
	"compress/lzw"
	"io"

	"github.com/goxmpp/goxmpp/extensions/features/compression"
)

func init() {
	compression.Methods = append(compression.Methods, compressor{compression.BaseCompressor{MethodName: "lzw"}})
}

type State struct {
	Level int
}

type compressor struct {
	compression.BaseCompressor
}

func (c compressor) GetReader(r io.Reader) (io.ReadCloser, error) {
	return lzw.NewReader(r, lzw.LSB, 8), nil
}

func (c compressor) GetWriter(w io.Writer) io.WriteCloser {
	return lzw.NewWriter(w, lzw.LSB, 8)
}
