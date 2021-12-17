package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordWithNoBitFieldsError(t *testing.T) {
	const (
		formatName   = "Format"
		functionName = "Marshal"
		wordName     = "Word"

		errorMessage = "" +
			"A word-struct should have exported fields " +
			"representing bit fields. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has no bit fields."
	)

	var (
		e error
	)

	e = NewWordWithNoBitFieldsError(functionName, formatName, wordName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
