package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldTypeUnsupportedError(t *testing.T) {
	const (
		fieldName  = "fieldName"
		formatName = "formatName"
		typeName   = "typeName"
		wordName   = "wordName"

		message = "Type typeName " +
			"of field fieldName " +
			"in word wordName " +
			"in format formatName " +
			"is not supported."
	)

	var (
		e error
	)

	e = NewFieldTypeUnsupportedError(formatName, wordName, fieldName, typeName)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
