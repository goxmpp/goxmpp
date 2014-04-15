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

type SwapRW func(source_rw io.ReadWriteCloser) (io.ReadWriteCloser, error)

func (self *Connection) UpdateRW(srw SwapRW) error {
	new_rw, err := srw(self.rw)
	if err == nil {
		self.rw = new_rw
		self.updateCoders()
	}
	return err
}

func (self *Connection) updateCoders() {
	self.streamEncoder = xml.NewEncoder(self.rw)
	self.streamDecoder = xml.NewDecoder(self.rw)
}

func (self *Connection) Close() error {
	return self.rw.Close()
}
