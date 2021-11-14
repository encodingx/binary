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
	// Return a slice of bytes encoding a Go struct
	// with properly defined field types and tags.

	// Even if an error occurs,
	// the byte slice returned should still be of the appropriate length.

	Unmarshal([]byte, interface{}) error
	// Parse a slice of bytes encoding a Go structs
	// with properly defined field types and tags, and
	// store the result in the structure pointed to by pointer.
}

type codec struct{}

func NewCodec() *codec {
	// Return a default implementation of interface Codec.

	return &codec{}
}

func (c codec) Marshal(structure interface{}) (bytes []byte, e error) {
	// Implement interface Codec.

	var (
		byteSlice []byte
		details   structFieldDetails
		field     bitfields.BitField
		i         int
		nFields   int
	)

	bytes = make([]byte, constants.WordLengthInBytes)

	nFields, details, e = getStructFieldDetails(structure)
	if e != nil {
		return
	}

	for i = 0; i < nFields; i++ {
		field, e = bitFieldFromStructFieldTag(details.Tags[i])
		if e != nil {
			return
		}

		switch details.Kinds[i] {
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
					details.Values[i].Uint(),
				),
			)
			if e != nil {
				return
			}

		case reflect.Bool:
			byteSlice, e = field.ByteSliceFromBool(
				details.Values[i].Bool(),
			)
			if e != nil {
				return
			}

		default:
			e = fmt.Errorf(unsupportedError,
				details.Kinds[i].String(),
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
		details     structFieldDetails
		field       bitfields.BitField
		i           int
		nFields     int
		valueBool   bool
		valueUint32 uint32
	)

	nFields, details, e = getStructFieldDetails(pointer)
	if e != nil {
		return
	}

	for i = 0; i < nFields; i++ {
		field, e = bitFieldFromStructFieldTag(details.Tags[i])
		if e != nil {
			return
		}

		switch details.Kinds[i] {
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

			details.Values[i].SetUint(
				uint64(valueUint32),
			)

		case reflect.Bool:
			valueBool, e = field.BoolFromByteSlice(bytes)
			if e != nil {
				return
			}

			details.Values[i].SetBool(valueBool)

		default:
			e = fmt.Errorf(unsupportedError,
				details.Kinds[i].String(),
			)

			return
		}
	}

	return
}

type structFieldDetails struct {
	Kinds  []reflect.Kind
	Tags   []reflect.StructTag
	Values []reflect.Value
}

func getStructFieldDetails(structure interface{}) (
	nFields int,
	details structFieldDetails,
	e error,
) {
	// Count the number of fields in a struct and
	// return a structFieldDetails containing three slices of that length
	// carrying the kinds, struct tags and values of those fields.

	const (
		nonStructError = "Type %s is not a struct."
	)

	var (
		structField reflect.StructField
		structType  reflect.Type
		structValue reflect.Value

		i int
	)

	structType = reflect.TypeOf(structure)

	if structType.Kind() != reflect.Struct {
		e = fmt.Errorf(nonStructError,
			structType.Name(),
		)

		return
	}

	structValue = reflect.ValueOf(structure)

	nFields = structType.NumField()

	details = structFieldDetails{
		Kinds:  make([]reflect.Kind, nFields),
		Tags:   make([]reflect.StructTag, nFields),
		Values: make([]reflect.Value, nFields),
	}

	for i = 0; i < nFields; i++ {
		structField = structType.Field(i)

		details.Kinds[i] = structField.Type.Kind()
		details.Tags[i] = structField.Tag
		details.Values[i] = structValue.Field(i)
	}

	return
}

func bitFieldFromStructFieldTag(tag reflect.StructTag) (
	field bitfields.BitField, e error,
) {
	// Return a bitfields.BitField given a struct field tag,
	// with the length and offset of the former set
	// to those indicated by the tag.

	var (
		length int
		offset int
	)

	length, offset, e = parseStructFieldTag(tag)
	if e != nil {
		return
	}

	field, e = bitfields.NewBitField(length, offset)
	if e != nil {
		return
	}

	return
}

func parseStructFieldTag(tag reflect.StructTag) (length, offset int, e error) {
	// Read values length and offset from a struct tag.

	const (
		key    = "bitfield"
		format = "%d,%d"
	)

	_, e = fmt.Sscanf(
		tag.Get(key),
		format,
		&length, &offset,
	)
	if e != nil {
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
