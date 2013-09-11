package decoder

import (
	"encoding/xml"
	"io"
)

type InnerXMLBuffer []byte

func (self *InnerXMLBuffer) Read(b []byte) (int, error) {
	if len(*self) == 0 {
		return 0, io.EOF
	}

	n := copy(b, *self)
	*self = (*self)[n:]
	return n, nil
}

type InnerDecoder struct {
	*xml.Decoder
	*InnerXMLBuffer
}

func NewInnerDecoder() *InnerDecoder {
	buffer := new(InnerXMLBuffer)
	return &InnerDecoder{
		Decoder:        xml.NewDecoder(buffer),
		InnerXMLBuffer: buffer,
	}
}

func (self *InnerDecoder) PutXML(b []byte) {
	*self.InnerXMLBuffer = b
}
