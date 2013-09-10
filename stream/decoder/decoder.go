package decoder

type InnerBuffer struct {
	buffer []byte
}

func NewInnerBuffer() *InnerBuffer {
	return &XMLBuffer{make([]byte, 0)}
}

func (self *InnerBuffer) Read(b []byte) (int, error) {
	if len(self.buffer) == 0 {
		return 0, io.EOF
	}

	n := copy(b, self.buffer)
	self.buffer = self.buffer[n:]
	return n, nil
}

func (self *InnerBuffer) PutXML(b []byte) {
	self.buffer = b
}

type InnerDecoder struct {
	*xml.Decoder
	*InnerBuffer
}

func NewInnerDecoder() *InnerDecoder {
	buffer := NewInnerBuffer()
	return &InnerDecoder{
		Decoder:     xml.NewDecoder(buffer),
		InnerBuffer: buffer,
	}
}
