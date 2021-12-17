package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitFieldOfUnsupportedTypeError(t *testing.T) {
	const (
		bitFieldName = "BitField"
		bitFieldType = "int"
		formatName   = "Format"
		functionName = "Marshal"
		wordName     = "Word"

		errorMessage = "" +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"of type uintN or bool. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has a bit field \"BitField\" " +
			"of unsupported type \"int\"."
	)

	var (
		e error
	)

	e = NewBitFieldOfUnsupportedTypeError(
		functionName,
		formatName,
		wordName,
		bitFieldName,
		bitFieldType,
	)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
