package compression_test

import (
	"bytes"
	"encoding/xml"
	"testing"
)

type connection struct {
	*bytes.Buffer
}

func NewConnection(b []byte) *connection {
	return &connection{Buffer: bytes.NewBuffer(b)}
}

func (c *connection) Close() error {
	return nil
}

// stream tags with following meaning
// [[<send>, <expect to get>], ...]
var stream = [][]string{
	{
		xml.Header + "<stream:stream xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' to='shakespeare.lit'>",
		xml.Header + "<stream:stream xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' to='shakespeare.lit'>",
	},
	{
		xml.Header + "<stream:stream xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' to='shakespeare.lit'>",
		xml.Header + "<stream:stream xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' to='shakespeare.lit'>",
	},
}

func TestCompression(t *testing.T) {

}
