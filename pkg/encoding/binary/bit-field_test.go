package binary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitField(t *testing.T) {
	const (
		lengthInvalid0 = 0
		lengthInvalid1 = 33
		lengthUint32   = 6
		lengthBool     = 1
		lengthValid0   = 1
		lengthValid1   = 32

		offsetCase0    = 26
		offsetCase1    = 13
		offsetCase2    = 0
		offsetCase3    = 27
		offsetInvalid0 = -1
		offsetInvalid1 = 32
		offsetValid0   = 0
		offsetValid1   = 31

		valueBool     = true
		valueOverflow = 1 << lengthUint32
		valueUint32   = 1<<lengthUint32 - 1 // as many ones as length in binary
	)

	var (
		byteSliceCase0  = []byte{0b11111100, 0b00000000, 0b00000000, 0b00000000}
		byteSliceCase1  = []byte{0b00000000, 0b00000111, 0b11100000, 0b00000000}
		byteSliceCase2  = []byte{0b00000000, 0b00000000, 0b00000000, 0b00111111}
		byteSliceCase3  = []byte{0b00001000, 0b00000000, 0b00000000, 0b00000000}
		byteSliceEmpty  = []byte{}
		byteSliceZeroes = []byte{0b00000000, 0b00000000, 0b00000000, 0b00000000}
	)

	var (
		bytes []byte
		e     error
		field BitField
		value interface{}
	)

	// Control cases
	// to make sure error in subsequent cases are not unintended.

	field, e = NewBitField(lengthValid0, offsetValid0)

	assert.Nil(t, e,
		"Function NewBitField should return a nil error "+
			"when it is passed a valid length argument.",
	)

	assert.NotNil(t, field,
		"Function NewBitField should return a non-nil BitField "+
			"when it is passed a valid length argument.",
	)

	// Error cases

	field, e = NewBitField(lengthInvalid0, offsetValid0)

	assert.NotNil(t, e,
		"Function NewBitField should return a non-nil error "+
			"when it is passed an invalid length argument.",
	)

	assert.Nil(t, field,
		"Function NewBitField should return a nil BitField "+
			"when it is passed an invalid length argument.",
	)

	field, e = NewBitField(lengthInvalid1, offsetValid1)

	assert.NotNil(t, e,
		"Function NewBitField should return a non-nil error "+
			"when it is passed an invalid length argument.",
	)

	assert.Nil(t, field,
		"Function NewBitField should return a nil BitField "+
			"when it is passed an invalid length argument.",
	)

	field, e = NewBitField(lengthValid1, offsetValid1)
	// valid length and offset but invalid combination

	assert.NotNil(t, e,
		"Function NewBitField should return a non-nil error "+
			"when it is passed an invalid length-offset combination.",
	)

	assert.Nil(t, field,
		"Function NewBitField should return a nil BitField "+
			"when it is passed an invalid length-offset combination.",
	)

	field, e = NewBitField(lengthUint32, offsetValid0)

	bytes, e = field.ByteSliceFromUint32(valueOverflow)

	assert.NotNil(t, e,
		"Method BitField.ByteSliceFromUint32 should "+
			"return a non-nil error "+
			"when it is passed a value that overflows the field.",
	)

	assert.Equal(t,
		byteSliceZeroes, bytes,
		"Method BitField.ByteSliceFromUint32 should "+
			"return a byte slice filled with zeroes "+
			"when it is passed a value that overflows the field.",
	)

	_, e = field.Uint32FromByteSlice(byteSliceEmpty)

	assert.NotNil(t, e,
		"Method BitField.Uint32FromByteSlice should "+
			"return a non-nil error "+
			"when it is passed a byte slice that is shorter than 4 bytes.",
	)

	_, e = field.BoolFromByteSlice(byteSliceEmpty)

	assert.NotNil(t, e,
		"Method BitField.BoolFromByteSlice should "+
			"return a non-nil error "+
			"when it is passed a byte slice that is shorter than 4 bytes.",
	)

	// Test Case 0

	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 <index
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |1 1 1 1 1 1                                                    |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//  3   2                   1                   0
	//  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 <offset

	field, e = NewBitField(lengthUint32, offsetCase0)
	if e != nil {
		t.Error(e)
	}

	bytes, e = field.ByteSliceFromUint32(valueUint32)
	if e != nil {
		t.Error(e)
	}

	assert.Equal(t,
		byteSliceCase0, bytes,
		"Byte slice returned by method BitField.ByteSliceFromUint32 "+
			"should match the expected.",
	)

	value, e = field.Uint32FromByteSlice(byteSliceCase0)
	if e != nil {
		t.Error(e)
	}

	assert.EqualValues(t,
		valueUint32, value,
		"Value returned by method BitField.Uint32FromByteSlice "+
			"should match the expected.",
	)

	// Test Case 1

	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 <index
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                          1 1 1 1 1 1                          |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//  3   2                   1                   0
	//  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 <offset

	field, e = NewBitField(lengthUint32, offsetCase1)
	if e != nil {
		t.Error(e)
	}

	bytes, e = field.ByteSliceFromUint32(valueUint32)
	if e != nil {
		t.Error(e)
	}

	assert.Equal(t,
		byteSliceCase1, bytes,
		"Byte slice returned by method BitField.ByteSliceFromUint32 "+
			"should match the expected.",
	)

	value, e = field.Uint32FromByteSlice(byteSliceCase1)
	if e != nil {
		t.Error(e)
	}

	assert.EqualValues(t,
		valueUint32, value,
		"Value returned by method BitField.Uint32FromByteSlice "+
			"should match the expected.",
	)

	// Test Case 2

	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 <index
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                                                    1 1 1 1 1 1|
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//  3   2                   1                   0
	//  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 <offset

	field, e = NewBitField(lengthUint32, offsetCase2)
	if e != nil {
		t.Error(e)
	}

	bytes, e = field.ByteSliceFromUint32(valueUint32)
	if e != nil {
		t.Error(e)
	}

	assert.Equal(t,
		byteSliceCase2, bytes,
		"Byte slice returned by method BitField.ByteSliceFromUint32 "+
			"should match the expected.",
	)

	value, e = field.Uint32FromByteSlice(byteSliceCase2)
	if e != nil {
		t.Error(e)
	}

	assert.EqualValues(t,
		valueUint32, value,
		"Value returned by method BitField.Uint32FromByteSlice "+
			"should match the expected.",
	)

	// Test Case 3

	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 <index
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |              1                                                |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//  3   2                   1                   0
	//  1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 <offset

	field, e = NewBitField(lengthBool, offsetCase3)
	if e != nil {
		t.Error(e)
	}

	bytes, e = field.ByteSliceFromBool(valueBool)
	if e != nil {
		t.Error(e)
	}

	assert.Equal(t,
		byteSliceCase3, bytes,
		"Byte slice returned by method BitField.ByteSliceFromBool "+
			"should match the expected.",
	)

	value, e = field.BoolFromByteSlice(byteSliceCase3)
	if e != nil {
		t.Error(e)
	}

	assert.EqualValues(t,
		valueBool, value,
		"Value returned by method BitField.BoolFromByteSlice "+
			"should match the expected.",
	)
}
