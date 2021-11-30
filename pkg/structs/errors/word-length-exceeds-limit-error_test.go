package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordLengthExceedsLimitError(t *testing.T) {
	const (
		formatName      = "formatName"
		wordName        = "wordName"
		wordLength      = 72
		wordLengthLimit = 64

		message = "Length 72 bits " +
			"of word wordName " +
			"in format formatName " +
			"exceeds the word length limit of 64 bits."
	)

	var (
		e error
	)

	e = NewWordLengthExceedsLimitError(
		formatName, wordName,
		wordLength, wordLengthLimit,
	)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
