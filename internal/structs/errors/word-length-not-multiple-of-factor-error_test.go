package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordLengthNotMultipleOfFactorError(t *testing.T) {
	const (
		factor     = 8
		formatName = "formatName"
		wordLength = 13
		wordName   = "wordName"

		message = "Length 13 bits " +
			"of word wordName " +
			"in format formatName " +
			"is not a multiple of 8."
	)

	var (
		e error
	)

	e = NewWordLengthNotMultipleOfFactorError(
		formatName, wordName,
		wordLength, factor,
	)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
