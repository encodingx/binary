package bitfields

import (
	"encoding/binary"
	"fmt"
)

const (
	Sub32BitFieldMaxSizeBits  = Sub32BitFieldMaxSizeBytes * 8
	Sub32BitFieldMaxSizeBytes = 4
	Sub32BitFieldMaxValue     = (1 << Sub32BitFieldMaxSizeBits) - 1
)

type Sub32BitField interface {
	// A bit field of length less than or equal to 32 bits.

	ByteSliceFromUint32(uint32) ([]byte, error)
	Uint32FromByteSlice([]byte) (uint32, error)

	ByteSliceFromBool(bool) ([]byte, error)
	BoolFromByteSlice([]byte) (bool, error)
}

type sub32BitField struct {
	length int
	// The bit-length of the field, not to be confused with that of its value.

	offset int
	// The number of places the bit field should be shifted left
	// from the rightmost section of a 32-bit sequence
	// for its position in that sequence to be appropriate.
}

func NewSub32BitField(length int, offset int) (field *sub32BitField, e error) {
	// Return a basic default implementation of interface Sub32BitField.

	field = &sub32BitField{
		length: length,
		offset: offset,
	}

	e = field.validateLengthAndOffset()
	if e != nil {
		field = nil

		return
	}

	return
}

func (f sub32BitField) ByteSliceFromUint32(input uint32) (
	output []byte, e error,
) {
	// Implement the Sub32BitField interface.

	return f.byteSlice(input)
}

func (f sub32BitField) Uint32FromByteSlice(input []byte) (
	output uint32, e error,
) {
	// Implement the Sub32BitField interface.

	return f.rawUint32(input)
}

func (f sub32BitField) ByteSliceFromBool(input bool) (output []byte, e error) {
	// Implement the Sub32BitField interface.

	// true is assumed to be equivalent to 1, and false to 0.
	// It is highly recommended that the length of boolean fields be set to 1
	// to minimise confusion.

	var (
		rawUint32 uint32
	)

	if input {
		rawUint32 = 1
	}

	return f.byteSlice(rawUint32)
}

func (f sub32BitField) BoolFromByteSlice(input []byte) (output bool, e error) {
	// Implement the Sub32BitField interface.

	// If the least-significant bit of the field is set to 1,
	// the value of the field is deemed to be true, otherwise false.
	// It is highly recommended that the length of boolean fields be set to 1
	// to minimise confusion.

	var (
		rawUint32 uint32
	)

	rawUint32, e = f.rawUint32(input)
	if e != nil {
		return
	}

	output = (rawUint32&1 == 1)

	return
}

func (f sub32BitField) byteSlice(rawUint32 uint32) (byteSlice []byte, e error) {
	// Given a 32-bit unsigned integer representing a value,
	// return a slice of four bytes representing a 32-bit sequence
	// containing the bit field in its appropriate position,
	// flanked by leading and trailing zeroes if applicable.

	// Even if an error occurs,
	// the byte slice returned should still be of the appropriate length.

	const (
		overflowError = "Unsigned integer %d overflows field of length %d."
	)

	byteSlice = make([]byte, Sub32BitFieldMaxSizeBytes)

	if rawUint32 >= (1 << f.length) {
		e = fmt.Errorf(overflowError, rawUint32, f.length)

		return
	}

	binary.BigEndian.PutUint32(byteSlice, rawUint32<<f.offset)

	return
}

func (f sub32BitField) rawUint32(byteSlice []byte) (rawUint32 uint32, e error) {
	// Given a slice of four bytes representing a 32-bit sequence,
	// return a 32-bit unsigned integer representing the value
	// contained by the bit field as defined by its length and offset.
	// Disregard bits in the sequence that fall outside the field.

	const (
		byteSliceLengthError = "Length of byte slice should be %d; got %d."
	)

	if len(byteSlice) != Sub32BitFieldMaxSizeBytes {
		e = fmt.Errorf(byteSliceLengthError,
			Sub32BitFieldMaxSizeBytes,
			len(byteSlice),
		)

		return
	}

	rawUint32 = binary.BigEndian.Uint32(byteSlice) & f.mask() >> f.offset

	return
}

func (f sub32BitField) mask() (mask uint32) {
	// Return the bit mask of the field
	// corresponding to its position in a 32-bit sequence.

	mask = Sub32BitFieldMaxValue >> (Sub32BitFieldMaxSizeBits - f.length) <<
		f.offset

	return
}

func (f sub32BitField) validateLengthAndOffset() (e error) {
	// Verify that the length and offset of a field
	// fall within appropriate ranges, and
	// that the length-offset combination
	// would not cause the field to overflow a 32-bit sequence.

	const (
		combinationError = "Combination of length %d and offset %d " +
			"would cause a sub-32-bit field to overflow a 32-bit sequence."
	)

	e = f.validateLength()
	if e != nil {
		return
	}

	e = f.validateOffset()
	if e != nil {
		return
	}

	if f.length+f.offset > Sub32BitFieldMaxSizeBits {
		e = fmt.Errorf(combinationError, f.length, f.offset)

		return
	}

	return
}

func (f sub32BitField) validateLength() (e error) {
	// Verify that the length of a field falls within the appropriate range.

	const (
		maximumLength = Sub32BitFieldMaxSizeBits
		minimumLength = 1

		lengthOutOfRange = "" +
			"Length %d does not fall within the appropriate range [%d, %d] " +
			"for a sub-32-bit field."
	)

	if f.length > maximumLength || f.length < minimumLength {
		e = fmt.Errorf(lengthOutOfRange, f.length, minimumLength, maximumLength)

		return
	}

	return
}

func (f sub32BitField) validateOffset() (e error) {
	// Verify that the offset of a field falls within the appropriate range.

	const (
		maximumOffset = Sub32BitFieldMaxSizeBits - 1
		minimumOffset = 0

		offsetOutOfRange = "" +
			"Offset %d does not fall within the appropriate range [%d, %d] " +
			"for a sub-32-bit field."
	)

	if f.offset > maximumOffset || f.offset < minimumOffset {
		e = fmt.Errorf(offsetOutOfRange, f.offset, minimumOffset, maximumOffset)

		return
	}

	return
}
