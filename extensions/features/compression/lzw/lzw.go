package lzw

import (
	"compress/lzw"
	"io"

	"github.com/dotdoom/goxmpp/extensions/features/compression"
	"github.com/dotdoom/goxmpp/stream"
)

func init() {
	compression.AddMethod("lzw", func(stream *stream.Stream) (compression.Compressor, error) {
		var state *State
		if err := stream.State.Get(&state); err != nil {
			return nil, err
		}
		return &compressor{}, nil
	})
}

type State struct {
	Level int
}

type compressor struct{}

func (c *compressor) GetReader(r io.Reader) (io.ReadCloser, error) {
	return lzw.NewReader(r, lzw.LSB, 8), nil
}

func (c *compressor) GetWriter(w io.Writer) io.WriteCloser {
	return lzw.NewWriter(w, lzw.LSB, 8)
}
