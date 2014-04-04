package stream

import (
	"encoding/xml"
	"io"
)

type Connection struct {
	rw            io.ReadWriteCloser
	streamEncoder *xml.Encoder
	streamDecoder *xml.Decoder
}

func (self *Connection) SetRW(rw io.ReadWriteCloser) {
	self.rw = rw
	self.updateCoders()
}

type SwapRW func(source_rw io.ReadWriteCloser) io.ReadWriteCloser

func (self *Connection) UpdateRW(srw SwapRW) {
	self.rw = srw(self.rw)
	self.updateCoders()
}

func (self *Connection) updateCoders() {
	self.streamEncoder = xml.NewEncoder(self.rw)
	self.streamDecoder = xml.NewDecoder(self.rw)
}
