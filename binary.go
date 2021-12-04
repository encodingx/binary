package binary

import (
	stdlib "encoding/binary"
	"io"

	"github.com/encodingx/binary/internal/codecs"
)

// Original features

var (
	defaultCodec = codecs.NewV1Codec()
)

func Marshal(pointer interface{}) ([]byte, error) {
	return defaultCodec.Marshal(pointer)
}

func Unmarshal(bytes []byte, pointer interface{}) error {
	return defaultCodec.Unmarshal(bytes, pointer)
}

// Standard library features

const (
	MaxVarintLen16 = stdlib.MaxVarintLen16
	MaxVarintLen32 = stdlib.MaxVarintLen32
	MaxVarintLen64 = stdlib.MaxVarintLen64
)

var (
	BigEndian    = stdlib.BigEndian
	LittleEndian = stdlib.LittleEndian
)

type ByteOrder interface {
	Uint16([]byte) uint16
	Uint32([]byte) uint32
	Uint64([]byte) uint64
	PutUint16([]byte, uint16)
	PutUint32([]byte, uint32)
	PutUint64([]byte, uint64)
	String() string
}

func PutUvarint(buf []byte, x uint64) int {
	return stdlib.PutUvarint(buf, x)
}

func PutVarint(buf []byte, x int64) int {
	return stdlib.PutVarint(buf, x)
}

func Read(r io.Reader, order ByteOrder, data interface{}) error {
	return stdlib.Read(r, order, data)
}

func ReadUvarint(r io.ByteReader) (uint64, error) {
	return stdlib.ReadUvarint(r)
}

func ReadVarint(r io.ByteReader) (int64, error) {
	return stdlib.ReadVarint(r)
}

func Size(v interface{}) int {
	return stdlib.Size(v)
}

func Uvarint(buf []byte) (uint64, int) {
	return stdlib.Uvarint(buf)
}

func Varint(buf []byte) (int64, int) {
	return stdlib.Varint(buf)
}

func Write(w io.Writer, order ByteOrder, data interface{}) error {
	return stdlib.Write(w, order, data)
}
