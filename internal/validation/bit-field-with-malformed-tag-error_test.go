package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitFieldWithMalformedTagError(t *testing.T) {
	const (
		bitFieldName = "BitField"
		formatName   = "Format"
		functionName = "Marshal"
		wordName     = "Word"

		errorMessage = "" +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"tagged with a key \"bitfield\" and a value " +
			"indicating the length of the bit field in number of bits " +
			"(e.g. `bitfield:\"1\"`). " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has a bit field \"BitField\" " +
			"with a malformed struct tag."
	)

	var (
		e error
	)

	e = NewBitFieldWithMalformedTagError(
		functionName,
		formatName,
		wordName,
		bitFieldName,
	)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
