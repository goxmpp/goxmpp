package connection

import (
	"encoding/xml"
	"io"
)

type State map[string]interface{}

type Connection struct {
	State         State
	RW            io.ReadWriter
	StreamEncoder *xml.Encoder
	StreamDecoder *xml.Decoder
}

func NewConnection(rw io.ReadWriter) *Connection {
	connection := &Connection{}
	connection.SetIO(rw)
	return connection
}

func (self *Connection) SetIO(rw io.ReadWriter) {
	self.RW = rw
	self.StreamEncoder = xml.NewEncoder(rw)
	self.StreamDecoder = xml.NewDecoder(rw)
}
