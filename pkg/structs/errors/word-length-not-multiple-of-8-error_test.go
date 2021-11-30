package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordLengthNotMultipleOf8Error(t *testing.T) {
	const (
		formatName = "formatName"
		wordName   = "wordName"
		wordLength = 13

		message = "Length 13 bits " +
			"of word wordName " +
			"in format formatName " +
			"is not a multiple of eight."
	)

	var (
		e error
	)

	e = NewWordLengthNotMultipleOf8Error(formatName, wordName, wordLength)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
