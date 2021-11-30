package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldStructTagMalformedError(t *testing.T) {
	const (
		fieldName  = "fieldName"
		formatName = "formatName"
		wordName   = "wordName"

		message = "Field fieldName " +
			"in word wordName " +
			"in format formatName " +
			"has a malformed struct tag."
	)

	var (
		e error
	)

	e = NewFieldStructTagMalformedError(formatName, wordName, fieldName)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
