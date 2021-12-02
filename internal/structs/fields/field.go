package fields

import (
	"fmt"
	"reflect"

	"github.com/joel-ling/go-bitfields/internal/structs/errors"
)

type Field interface {
	// A Field represents a sequence of bits in a binary message or file format.

	LengthInBits() uint
	// Return the number of bits in a Field.

	OffsetInBits() uint
	// Return the number of bits to the right of the rightmost bit of a Field
	// in a Word.

	Kind() reflect.Kind
	// Return the reflect.Kind of the Field.
}

type field struct {
	lengthInBits uint
	offsetInBits uint
	kind         reflect.Kind
}

func NewFieldFromStructField(structField reflect.StructField,
	wordName, formatName string,
	sumFieldLengthsBefore uint,
) (
	f *field, sumFieldLengthsAfter uint, e error,
) {
	// Return a default implementation of the interface Field
	// and a cumulative sum of the length of some collection of Fields,
	// given a reflection of a struct field representing a Field in a Word,
	// the names of the Format and Word containing that Field,
	// and the previous cumulative sum of the length of those Fields.

	const (
		structFieldTagKey         = "bitfield"
		structFieldTagValueFormat = "%d,%d"

		sizeBool   = 1
		sizeUint   = 64
		sizeUint16 = 16
		sizeUint32 = 32
		sizeUint64 = 64
		sizeUint8  = 8
	)

	var (
		structFieldTagKeyOK bool
		structFieldTagValue string
		structFieldKindSize uint
	)

	f = &field{
		kind: structField.Type.Kind(),
	}

	switch f.kind {
	case reflect.Uint:
		structFieldKindSize = sizeUint

	case reflect.Uint8:
		structFieldKindSize = sizeUint8

	case reflect.Uint16:
		structFieldKindSize = sizeUint16

	case reflect.Uint32:
		structFieldKindSize = sizeUint32

	case reflect.Uint64:
		structFieldKindSize = sizeUint64

	case reflect.Bool:
		structFieldKindSize = sizeBool

	default:
		e = errors.NewFieldTypeUnsupportedError(
			formatName,
			wordName,
			structField.Name,
			f.kind.String(),
		)

		return
	}

	if len(structField.Tag) == 0 {
		e = errors.NewFieldMissingStructTagError(
			formatName,
			wordName,
			structField.Name,
		)

		return
	}

	structFieldTagValue, structFieldTagKeyOK = structField.Tag.Lookup(
		structFieldTagKey,
	)

	_, e = fmt.Sscanf(
		structFieldTagValue,
		structFieldTagValueFormat,
		&f.lengthInBits,
		&f.offsetInBits,
	)
	if e != nil || !structFieldTagKeyOK {
		e = errors.NewFieldStructTagMalformedError(
			formatName,
			wordName,
			structField.Name,
		)

		return
	}

	if f.lengthInBits > structFieldKindSize {
		e = errors.NewFieldLengthOverflowsTypeError(
			formatName,
			wordName,
			structField.Name,
			structField.Type.Kind().String(),
			structFieldKindSize,
			f.lengthInBits,
		)

		return
	}

	if f.offsetInBits != sumFieldLengthsBefore {
		e = errors.NewFieldHasGapOrOverlapToRightError(
			formatName,
			wordName,
			structField.Name,
		)

		return
	}

	sumFieldLengthsAfter = sumFieldLengthsBefore + f.lengthInBits

	return
}

func (f *field) LengthInBits() uint {
	// Implement interface Field.

	return f.lengthInBits
}

func (f *field) OffsetInBits() uint {
	// Implement interface Field.

	return f.offsetInBits
}

func (f *field) Kind() reflect.Kind {
	// Implement interface Field.

	return f.kind
}
