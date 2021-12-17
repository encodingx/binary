package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	bitFieldName = "BitField"
)

func TestBitFieldOfUnsupportedTypeError(t *testing.T) {
	const (
		bitFieldType = "int"

		errorMessage = "" +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"of type uintN or bool. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has a bit field \"BitField\" " +
			"of unsupported type \"int\"."
	)

	var (
		e BitFieldError
	)

	e = NewBitFieldOfUnsupportedTypeError(bitFieldType)

	e.SetFunctionName(functionName)

	e.SetFormatName(formatName)

	e.SetWordName(wordName)

	e.SetBitFieldName(bitFieldName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}

func TestBitFieldOfLengthOverflowingTypeError(t *testing.T) {
	const (
		bitFieldLength = 32
		bitFieldType   = "uint8"

		errorMessage = "" +
			"The number of unique values a bit field can contain " +
			"must not exceed the size of its type. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has a bit field \"BitField\" " +
			"of length 32 exceeding the size of type \"uint8\"."
	)

	var (
		e BitFieldError
	)

	e = NewBitFieldOfLengthOverflowingTypeError(bitFieldLength, bitFieldType)

	e.SetFunctionName(functionName)

	e.SetFormatName(formatName)

	e.SetWordName(wordName)

	e.SetBitFieldName(bitFieldName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}

func TestBitFieldWithMalformedTagError(t *testing.T) {
	const (
		errorMessage = "" +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"tagged with a key \"bitfield\" and a value " +
			"indicating the length of the bit field in number of bits " +
			"(e.g. `bitfield:\"1\"`). " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has a bit field \"BitField\" " +
			"with a malformed struct tag."
	)

	var (
		e BitFieldError
	)

	e = NewBitFieldWithMalformedTagError()

	e.SetFunctionName(functionName)

	e.SetFormatName(formatName)

	e.SetWordName(wordName)

	e.SetBitFieldName(bitFieldName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}

func TestBitFieldWithNoStructTagError(t *testing.T) {
	const (
		errorMessage = "" +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"tagged with a key \"bitfield\" and a value " +
			"indicating the length of the bit field in number of bits " +
			"(e.g. `bitfield:\"1\"`). " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has a bit field \"BitField\" " +
			"with no struct tag."
	)

	var (
		e BitFieldError
	)

	e = NewBitFieldWithNoStructTagError()

	e.SetFunctionName(functionName)

	e.SetFormatName(formatName)

	e.SetWordName(wordName)

	e.SetBitFieldName(bitFieldName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
