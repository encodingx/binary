package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordNotStructError(t *testing.T) {
	const (
		formatName   = "Format"
		functionName = "Marshal"
		wordName     = "Word"

		errorMessage = "" +
			"A format-struct should nest exported word-structs. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"that is not a struct."
	)

	var (
		e error
	)

	e = NewWordNotStructError(functionName, formatName, wordName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
