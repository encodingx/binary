package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitFieldOfLengthOverflowingTypeError(t *testing.T) {
	const (
		bitFieldName = "BitField"
		bitFieldType = "uint8"
		formatName   = "Format"
		functionName = "Marshal"
		wordName     = "Word"

		bitFieldLength = 32

		errorMessage = "" +
			"The number of unique values a bit field can contain " +
			"must not exceed the size of its type. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has a bit field \"BitField\" " +
			"of length 32 exceeding the size of type \"uint8\"."
	)

	var (
		e error
	)

	e = NewBitFieldOfLengthOverflowingTypeError(
		functionName,
		formatName,
		wordName,
		bitFieldName,
		bitFieldLength,
		bitFieldType,
	)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
