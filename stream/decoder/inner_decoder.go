package decoder

import (
	"encoding/xml"
	"io"
)

const TERMINATOR = "nonexisting_node_this_is_a_hack"

type InnerXMLBuffer []byte

func (self *InnerXMLBuffer) ReadByte() (byte, error) {
	var b byte
	var err error
	if len(*self) > 0 {
		b = (*self)[0]
		*self = (*self)[1:]
	} else {
		err = io.EOF
	}

	return b, err
}

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
	// We need this fake tag to trick xml.Decoder and exit the parsing loop once inner XML buffer is empty.
	// We check if buffer is empty after every tag parsed. Without this tag we can have "chardata" before
	// outer closing tag and xml.Decoder.Token will never return till EOF,
	// but after EOF it won't read any more data from buffer even if we provide any
	*self.InnerXMLBuffer = append(*self.InnerXMLBuffer, "<"+TERMINATOR+"/>"...)
}

func (self *InnerDecoder) IsEmpty() bool {
	return len(*self.InnerXMLBuffer) == 0
}
