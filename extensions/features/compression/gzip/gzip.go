package gzip

import (
	"compress/gzip"
	"io"

	"github.com/dotdoom/goxmpp/extensions/features/compression"
	"github.com/dotdoom/goxmpp/stream"
)

func init() {
	compression.AddMethod("gzip", func(stream *stream.Stream) (compression.Compressor, error) {
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
	return gzip.NewReader(r)
}

func (c *compressor) GetWriter(w io.Writer) io.WriteCloser {
	return gzip.NewWriter(w)
}
