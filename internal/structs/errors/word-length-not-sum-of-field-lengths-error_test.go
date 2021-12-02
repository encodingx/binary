package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordLengthNotSumOfFieldLengthsError(t *testing.T) {
	const (
		formatName        = "formatName"
		wordName          = "wordName"
		wordLength        = 16
		sumOfFieldLengths = 24

		message = "Length 16 bits declared in struct tag " +
			"of word wordName " +
			"in format formatName " +
			"is not equal to the sum of the lengths of its fields, 24 bits."
	)

	var (
		e error
	)

	e = NewWordLengthNotSumOfFieldLengthsError(
		formatName, wordName,
		wordLength, sumOfFieldLengths,
	)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
