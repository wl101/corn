package internal

import "encoding/binary"

var metaHeaderSize int64 = 8

type meta struct {
	keySize uint32
	valSize uint32
	key     []byte
	val     []byte
}

func (m *meta) Encode() []byte {
	buf := make([]byte, int(metaHeaderSize)+len(m.key)+len(m.val))
	binary.BigEndian.PutUint32(buf[:4], m.keySize)
	binary.BigEndian.PutUint32(buf[4:8], m.valSize)
	n := copy(buf[metaHeaderSize:], m.key)
	copy(buf[int(metaHeaderSize)+n:], m.val)
	return buf
}

func (m *meta) Decode(buf []byte) {
	if m.keySize == 0 {
		m.keySize = binary.BigEndian.Uint32(buf[:4])
		m.valSize = binary.BigEndian.Uint32(buf[4:])
	} else {
		m.key = buf[:m.keySize]
		m.val = buf[m.keySize:]
	}
}
