package stream

import (
	"encoding/xml"
	"io"
)

type Connection struct {
	rw            io.ReadWriter
	streamEncoder *xml.Encoder
	streamDecoder *xml.Decoder
}

func (self *Connection) SetRW(rw io.ReadWriter) io.ReadWriter {
	previous_rw := self.rw
	self.rw = rw
	self.streamEncoder = xml.NewEncoder(rw)
	self.streamDecoder = xml.NewDecoder(rw)
	return previous_rw
}
