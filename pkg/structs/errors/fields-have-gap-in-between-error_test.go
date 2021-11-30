package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldsHaveGapInBetweenError(t *testing.T) {
	const (
		formatName     = "formatName"
		wordName       = "wordName"
		leftFieldName  = "leftFieldName"
		rightFieldName = "rightFieldName"

		message = "There is a gap " +
			"between fields leftFieldName " +
			"and rightFieldName " +
			"in word wordName " +
			"in format formatName."
	)

	var (
		e error
	)

	e = NewFieldsHaveGapInBetweenError(
		formatName, wordName,
		leftFieldName, rightFieldName,
	)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
