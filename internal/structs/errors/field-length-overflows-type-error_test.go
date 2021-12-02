package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldLengthOverflowsTypeError(t *testing.T) {
	const (
		fieldName   = "fieldName"
		formatName  = "formatName"
		typeName    = "typeName"
		wordName    = "wordName"
		fieldLength = 9
		typeLength  = 8

		message = "Length 9 bits " +
			"of field fieldName " +
			"in word wordName " +
			"in format formatName " +
			"overflows type typeName " +
			"of length 8 bits."
	)

	var (
		e error
	)

	e = NewFieldLengthOverflowsTypeError(
		formatName, wordName, fieldName, typeName,
		typeLength, fieldLength,
	)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
