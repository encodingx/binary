package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	formatName = "Format"
)

func TestFormatWithNoWordsError(t *testing.T) {
	const (
		errorMessage = "" +
			"A format-struct should nest exported word-structs. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"that has no words."
	)

	var (
		e FormatError
	)

	e = NewFormatWithNoWordsError()

	e.SetFunctionName(functionName)

	e.SetFormatName(formatName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}

func TestLengthOfByteSliceNotEqualToFormatLengthError(t *testing.T) {
	const (
		byteSliceLength     = 1
		formatLengthInBytes = 4

		errorMessage = "" +
			"A byte slice into which a format-struct would be unmarshalled " +
			"should be of length equal to the sum of lengths of words " +
			"in the format represented by the struct. " +
			"Argument to Unmarshal points to a format-struct \"Format\" " +
			"of length 4 byte(s) " +
			"not equal to the length of the byte slice, 1 byte(s)."
	)

	var (
		e FormatError
	)

	e = NewLengthOfByteSliceNotEqualToFormatLengthError(
		formatLengthInBytes,
		byteSliceLength,
	)

	e.SetFormatName(formatName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
