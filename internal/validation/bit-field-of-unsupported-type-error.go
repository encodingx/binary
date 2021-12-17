package validation

import (
	"fmt"
)

func NewBitFieldOfUnsupportedTypeError(
	functionName, formatName, wordName, bitFieldName, bitFieldType string,
) (
	e error,
) {
	const (
		format = "" +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"of type uintN or bool. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"that has a bit field \"%s\" " +
			"of unsupported type \"%s\"."
	)

	e = fmt.Errorf(format,
		functionName, formatName, wordName, bitFieldName, bitFieldType,
	)

	return
}
