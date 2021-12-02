package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatHasNoWordsError(t *testing.T) {
	const (
		formatName = "formatName"

		message = "Format formatName has no words in it."
	)

	var (
		e error
	)

	e = NewFormatHasNoWordsError(formatName)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
