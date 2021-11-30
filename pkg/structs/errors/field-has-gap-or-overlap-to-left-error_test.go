package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldHasGapOrOverlapToLeftError(t *testing.T) {
	const (
		fieldName  = "fieldName"
		formatName = "formatName"
		wordName   = "wordName"

		message = "There is a gap or overlap " +
			"to the left of field fieldName " +
			"in word wordName " +
			"in format formatName."
	)

	var (
		e error
	)

	e = NewFieldHasGapOrOverlapToLeftError(formatName, wordName, fieldName)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
