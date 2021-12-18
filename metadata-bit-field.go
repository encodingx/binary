package binary

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/encodingx/binary/internal/validation"
)

type bitFieldMetadata struct {
	length uint
	offset uint64
	kind   reflect.Kind
}

func newBitFieldMetadataFromStructFieldReflection(
	reflection reflect.StructField,
) (
	bitField bitFieldMetadata, e error,
) {
	const (
		tagKey         = "bitfield"
		tagValueFormat = "%d"
	)

	var (
		bitFieldLengthCap uint
	)

	defer func() {
		if e != nil {
			e.(validation.BitFieldError).SetBitFieldName(reflection.Name)
		}
	}()

	switch reflection.Type.Kind() {
	case reflect.Uint:
		fallthrough

	case reflect.Uint64:
		bitFieldLengthCap = 64

	case reflect.Uint32:
		bitFieldLengthCap = 32

	case reflect.Uint16:
		bitFieldLengthCap = 16

	case reflect.Uint8:
		bitFieldLengthCap = 8

	case reflect.Bool:
		bitFieldLengthCap = 1

	default:
		e = validation.NewBitFieldOfUnsupportedTypeError(
			reflection.Type.String(),
		)

		return
	}

	bitField = bitFieldMetadata{
		kind: reflection.Type.Kind(),
	}

	if len(reflection.Tag) == 0 {
		e = validation.NewBitFieldWithNoStructTagError()

		return
	}

	_, e = fmt.Sscanf(
		reflection.Tag.Get(tagKey),
		tagValueFormat,
		&bitField.length,
	)
	if e != nil {
		e = validation.NewBitFieldWithMalformedTagError()

		return
	}

	if bitField.length > bitFieldLengthCap {
		e = validation.NewBitFieldOfLengthOverflowingTypeError(
			bitField.length,
			reflection.Type.String(),
		)
	}

	return
}

func (m bitFieldMetadata) marshal(reflection reflect.Value) (value uint64) {
	switch m.kind {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough

	case reflect.Uint:
		value = reflection.Uint()

	case reflect.Bool:
		if reflection.Bool() {
			value = 1
		}
	}

	value = value & (1<<m.length - 1) << m.offset // XXX: mask if overflowing

	return
}

func (m bitFieldMetadata) unmarshal(bytes []byte, reflection reflect.Value) {
	var (
		value uint64
	)

	value = binary.BigEndian.Uint64(bytes) >> m.offset & (1<<m.length - 1)

	switch m.kind {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough

	case reflect.Uint:
		reflection.SetUint(value)

	case reflect.Bool:
		switch value {
		case 1:
			reflection.SetBool(true)
		default:
			reflection.SetBool(false)
		}
	}

	return
}
