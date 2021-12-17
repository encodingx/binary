package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLengthOfByteSliceNotEqualToFormatLengthError(t *testing.T) {
	const (
		byteSliceLength = 8
		formatLength    = 32

		formatName = "Format"

		errorMessage = "" +
			"A byte slice into which a format-struct would be unmarshalled " +
			"should be of length equal to the sum of lengths of words " +
			"in the format represented by the struct. " +
			"Argument to Unmarshal points to a format-struct \"Format\" " +
			"of length 32 bits " +
			"not equal to the length of the byte slice, 8 bits."
	)

	var (
		e error
	)

	e = NewLengthOfByteSliceNotEqualToFormatLengthError(
		formatName,
		formatLength,
		byteSliceLength,
	)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
