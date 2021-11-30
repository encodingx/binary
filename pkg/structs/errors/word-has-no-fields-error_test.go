package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordHasNoFieldsError(t *testing.T) {
	const (
		formatName = "formatName"
		wordName   = "wordName"

		message = "Word wordName " +
			"in format formatName " +
			"has no fields in it."
	)

	var (
		e error
	)

	e = NewWordHasNoFieldsError(formatName, wordName)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
