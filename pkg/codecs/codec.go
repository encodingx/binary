package codecs

import (
	"fmt"
	"reflect"

	"github.com/joel-ling/go-bitfields/pkg/bitfields"
	"github.com/joel-ling/go-bitfields/pkg/constants"
)

const (
	unsupportedError = "Struct fields of kind %s are not yet supported."
)

type Codec interface {
	Marshal(interface{}) ([]byte, error)
	// Return a slice of four bytes containing a 32-bit "word"
	// encoding a Go struct with properly defined field types and tags.

	// Even if an error occurs,
	// the byte slice returned should still be of the appropriate length.

	Unmarshal([]byte, interface{}) error
	// Parse a slice of four bytes containing a 32-bit "word"
	// encoding a Go struct with properly defined field types and tags, and
	// store the result in the structure pointed to by pointer.
}

type codec struct {
	bitFieldCache map[reflect.StructTag]bitfields.BitField
}

func NewCodec() *codec {
	// Return a default implementation of interface Codec.

	return &codec{bitFieldCache: make(map[reflect.StructTag]bitfields.BitField)}
}

func (c codec) Marshal(pointer interface{}) (bytes []byte, e error) {
	// Implement interface Codec.

	var (
		byteSlice   []byte
		field       bitfields.BitField
		fieldType   reflect.StructField
		fieldValue  reflect.Value
		i           int
		nFields     int
		structType  reflect.Type
		structValue reflect.Value
	)

	bytes = make([]byte, constants.WordLengthInBytes)

	structType, structValue, e = reflectStruct(pointer)
	if e != nil {
		return
	}

	nFields = structType.NumField()

	for i = 0; i < nFields; i++ {
		fieldType = structType.Field(i)
		fieldValue = structValue.Field(i)

		field, e = c.bitFieldFromStructFieldTag(fieldType.Tag)
		if e != nil {
			return
		}

		switch fieldValue.Kind() {
		case reflect.Uint8:
			fallthrough

		case reflect.Uint16:
			fallthrough

		case reflect.Uint32:
			fallthrough

		case reflect.Uint64:
			fallthrough

		case reflect.Uint:
			byteSlice, e = field.ByteSliceFromUint32(
				uint32(
					fieldValue.Uint(),
				),
			)
			if e != nil {
				return
			}

		case reflect.Bool:
			byteSlice, e = field.ByteSliceFromBool(
				fieldValue.Bool(),
			)
			if e != nil {
				return
			}

		default:
			e = fmt.Errorf(unsupportedError,
				fieldValue.Kind().String(),
			)

			return
		}

		bytes = byteSliceBitwiseOr(bytes, byteSlice)
	}

	return
}

func (c codec) Unmarshal(bytes []byte, pointer interface{}) (e error) {
	// Implement interface Codec.

	var (
		field       bitfields.BitField
		fieldType   reflect.StructField
		fieldValue  reflect.Value
		i           int
		nFields     int
		structType  reflect.Type
		structValue reflect.Value
		valueBool   bool
		valueUint32 uint32
	)

	structType, structValue, e = reflectStruct(pointer)
	if e != nil {
		return
	}

	nFields = structType.NumField()

	for i = 0; i < nFields; i++ {
		fieldType = structType.Field(i)
		fieldValue = structValue.Field(i)

		field, e = c.bitFieldFromStructFieldTag(fieldType.Tag)
		if e != nil {
			return
		}

		switch fieldValue.Kind() {
		case reflect.Uint8:
			fallthrough

		case reflect.Uint16:
			fallthrough

		case reflect.Uint32:
			fallthrough

		case reflect.Uint64:
			fallthrough

		case reflect.Uint:
			valueUint32, e = field.Uint32FromByteSlice(bytes)
			if e != nil {
				return
			}

			fieldValue.SetUint(
				uint64(valueUint32),
			)

		case reflect.Bool:
			valueBool, e = field.BoolFromByteSlice(bytes)
			if e != nil {
				return
			}

			fieldValue.SetBool(valueBool)

		default:
			e = fmt.Errorf(unsupportedError,
				fieldValue.Kind().String(),
			)

			return
		}
	}

	return
}

func (c codec) bitFieldFromStructFieldTag(tag reflect.StructTag) (
	field bitfields.BitField, e error,
) {
	// Return a bitfields.BitField given a struct field tag,
	// with the length and offset of the former set
	// to those indicated by the tag.

	var (
		cached bool
		length int
		offset int
	)

	field, cached = c.bitFieldCache[tag]
	if cached {
		return
	}

	length, offset, e = parseStructFieldTag(tag)
	if e != nil {
		return
	}

	field, e = bitfields.NewBitField(length, offset)
	if e != nil {
		return
	}

	c.bitFieldCache[tag] = field

	return
}

func reflectStruct(pointer interface{}) (
	goType reflect.Type, goValue reflect.Value, e error,
) {
	// Verify that pointer is indeed a pointer to a struct, and
	// return a reflect.Type and reflect.Value
	// representing the Go type and value of the struct.

	const (
		notPointerError = "Type %s is not a pointer."
		notStructError  = "Type %s is not a struct."
	)

	goType = reflect.TypeOf(pointer)
	goValue = reflect.ValueOf(pointer)

	if goType.Kind() != reflect.Ptr {
		e = fmt.Errorf(notPointerError,
			goType.String(),
		)

		return
	}

	goType = goType.Elem()
	goValue = goValue.Elem()

	if goType.Kind() != reflect.Struct {
		e = fmt.Errorf(notStructError,
			goType.String(),
		)

		return
	}

	return
}

func parseStructFieldTag(tag reflect.StructTag) (length, offset int, e error) {
	// Read values length and offset from a struct tag.

	const (
		key    = "bitfield"
		format = "%d,%d"

		tagError = "Could not parse struct tag `%s`: %w"
	)

	_, e = fmt.Sscanf(
		tag.Get(key),
		format,
		&length, &offset,
	)
	if e != nil {
		e = fmt.Errorf(tagError, tag, e)

		return
	}

	return
}

func byteSliceBitwiseOr(a []byte, b []byte) (c []byte) {
	// Perform the bitwise OR operation on byte slices a and b to derive c.
	// Length of slice c is the greater of the lengths of a and b.

	var (
		i      int
		length int
	)

	length = len(a)

	if len(b) > len(a) {
		length = len(b)
	}

	c = make([]byte, length)

	for i = 0; i < length; i++ {
		c[i] = a[i] | b[i]
	}

	return
}
