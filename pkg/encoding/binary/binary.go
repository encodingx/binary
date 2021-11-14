package binary

import (
	"encoding/binary"
	"io"
)

// Original features

var (
	defaultCodec = NewCodec()
)

func Marshal(structure interface{}) ([]byte, error) {
	return defaultCodec.Marshal(structure)
}

func Unmarshal(bytes []byte, pointer interface{}) error {
	return defaultCodec.Unmarshal(bytes, pointer)
}

// Standard library features

const (
	MaxVarintLen16 = binary.MaxVarintLen16
	MaxVarintLen32 = binary.MaxVarintLen32
	MaxVarintLen64 = binary.MaxVarintLen64
)

var (
	BigEndian    = binary.BigEndian
	LittleEndian = binary.LittleEndian
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
	return binary.PutUvarint(buf, x)
}

func PutVarint(buf []byte, x int64) int {
	return binary.PutVarint(buf, x)
}

func Read(r io.Reader, order ByteOrder, data interface{}) error {
	return binary.Read(r, order, data)
}

func ReadUvarint(r io.ByteReader) (uint64, error) {
	return binary.ReadUvarint(r)
}

func ReadVarint(r io.ByteReader) (int64, error) {
	return binary.ReadVarint(r)
}

func Size(v interface{}) int {
	return binary.Size(v)
}

func Uvarint(buf []byte) (uint64, int) {
	return binary.Uvarint(buf)
}

func Varint(buf []byte) (int64, int) {
	return binary.Varint(buf)
}

func Write(w io.Writer, order ByteOrder, data interface{}) error {
	return binary.Write(w, order, data)
}
