package decoder

import (
	"encoding/xml"
)

type InnerXMLBuffer struct {
	buffer []byte
}

func NewInnerXMLBuffer() *InnerXMLBuffer {
	return &InnerXMLBuffer{make([]byte, 0)}
}

func (self *InnerXMLBuffer) Read(b []byte) (int, error) {
	if len(self.buffer) == 0 {
		return 0, io.EOF
	}

	n := copy(b, self.buffer)
	self.buffer = self.buffer[n:]
	return n, nil
}

func (self *InnerXMLBuffer) PutXML(b []byte) {
	self.buffer = b
}

type InnerDecoder struct {
	*xml.Decoder
	*InnerXMLBuffer
}

func NewInnerDecoder() *InnerDecoder {
	buffer := NewInnerXMLBuffer()
	return &InnerDecoder{
		Decoder:        xml.NewDecoder(buffer),
		InnerXMLBuffer: buffer,
	}
}
