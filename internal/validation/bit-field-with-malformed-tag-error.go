package validation

import (
	"fmt"
)

func NewBitFieldWithMalformedTagError(
	functionName, formatName, wordName, bitFieldName string,
) (
	e error,
) {
	const (
		format = "" +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"tagged with a key \"bitfield\" and a value " +
			"indicating the length of the bit field in number of bits " +
			"(e.g. `bitfield:\"1\"`). " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"that has a bit field \"%s\" " +
			"with a malformed struct tag."
	)

	return fmt.Errorf(format, functionName, formatName, wordName, bitFieldName)
}
