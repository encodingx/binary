package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordIsNotStructError(t *testing.T) {
	const (
		formatName = "formatName"
		wordName   = "wordName"

		message = "Word wordName " +
			"in format formatName " +
			"is not a struct."
	)

	var (
		e error
	)

	e = NewWordIsNotStructError(formatName, wordName)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
