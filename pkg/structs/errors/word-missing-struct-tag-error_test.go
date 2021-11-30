package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordMissingStructTagError(t *testing.T) {
	const (
		formatName = "formatName"
		wordName   = "wordName"

		message = "Word wordName " +
			"in format formatName " +
			"is missing a struct tag."
	)

	var (
		e error
	)

	e = NewWordMissingStructTagError(formatName, wordName)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
