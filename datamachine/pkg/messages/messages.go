package messages

import (
	"bytes"
	"encoding/binary"
)

const (
	MessageSize int = 144
)

type Message struct {
	Time  int64
	Data  [16]byte
	Value int64
}

func (m *Message) UDPEncode() ([]byte, error) {
	b := make([]byte, MessageSize)
	if _, err := binary.Encode(b, binary.LittleEndian, *m); err != nil {
		return nil, err
	}
	return b, nil
}

func (m *Message) UDPDecode(b []byte) error {
	return binary.Read(bytes.NewReader(b), binary.LittleEndian, m)
}
