package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordOfLengthNotEqualToSumOfLengthsOfBitFieldsError(t *testing.T) {
	const (
		formatName   = "Format"
		functionName = "Marshal"
		wordName     = "Word"

		bitFieldLengthSum = 31
		wordLength        = 32

		errorMessage = "" +
			"The length of a word " +
			"should be equal to the sum of lengths of its bit fields. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"that has a word \"Word\" " +
			"of length 32 " +
			"not equal to the sum of the lengths of its bit fields, 31."
	)

	var (
		e error
	)

	e = NewWordOfLengthNotEqualToSumOfLengthsOfBitFieldsError(
		functionName, formatName, wordName,
		wordLength, bitFieldLengthSum,
	)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
