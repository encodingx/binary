package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordOfIncompatibleLengthError(t *testing.T) {
	const (
		formatName   = "Format"
		functionName = "Marshal"
		wordName     = "Word"

		wordLength = 36

		errorMessage = "" +
			"The length of a word should be a multiple of eight " +
			"in the range [8, 64]. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"that has a word \"Word\" " +
			"of length 36 not in {8, 16, 24, ... 64}."
	)

	var (
		e error
	)

	e = NewWordOfIncompatibleLengthError(functionName, formatName, wordName,
		wordLength,
	)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
