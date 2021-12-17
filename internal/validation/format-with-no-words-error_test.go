package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatWithNoWordsError(t *testing.T) {
	const (
		formatName   = "Format"
		functionName = "Marshal"

		errorMessage = "" +
			"A format-struct should nest exported word-structs. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"that has no words."
	)

	var (
		e error
	)

	e = NewFormatWithNoWordsError(functionName, formatName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
