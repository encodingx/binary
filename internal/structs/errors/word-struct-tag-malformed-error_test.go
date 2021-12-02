package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordStructTagMalformedError(t *testing.T) {
	const (
		formatName = "formatName"
		wordName   = "wordName"

		message = "Word wordName " +
			"in format formatName " +
			"has a malformed struct tag."
	)

	var (
		e error
	)

	e = NewWordStructTagMalformedError(formatName, wordName)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
