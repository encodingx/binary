package binary

import (
	stdlib "encoding/binary"
	"fmt"
	"io"

	"github.com/encodingx/binary/internal/validation"
)

const (
	wordLengthUpperLimitBytes = 8
)

var (
	defaultCodec = newCodec()
)

func Marshal(iface interface{}) (bytes []byte, e error) {
	const (
		functionName = "Marshal"
	)

	var (
		operation codecOperation
	)

	defer func() {
		const (
			marshalError = "Marshal error: %w"
		)

		if e != nil {
			e.(validation.FunctionError).SetFunctionName(functionName)

			e = fmt.Errorf(marshalError, e)
		}

		return
	}()

	operation, e = defaultCodec.newOperation(iface)
	if e != nil {
		return
	}

	bytes, e = operation.marshal()
	if e != nil {
		return
	}

	return
}

func Unmarshal(bytes []byte, iface interface{}) (e error) {
	const (
		functionName = "Unmarshal"
	)

	var (
		operation codecOperation
	)

	defer func() {
		const (
			unmarshalError = "Unmarshal error: %w"
		)

		if e != nil {
			e.(validation.FunctionError).SetFunctionName(functionName)

			e = fmt.Errorf(unmarshalError, e)
		}

		return
	}()

	operation, e = defaultCodec.newOperation(iface)
	if e != nil {
		return
	}

	e = operation.unmarshal(bytes)
	if e != nil {
		return
	}

	return
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
