package validation

import (
	"fmt"
)

func NewBitFieldOfLengthOverflowingTypeError(
	functionName, formatName, wordName, bitFieldName string,
	bitFieldLength uint, bitFieldType string,
) (
	e error,
) {
	const (
		format = "" +
			"The number of unique values a bit field can contain " +
			"must not exceed the size of its type. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"that has a bit field \"%s\" " +
			"of length %d exceeding the size of type \"%s\"."
	)

	e = fmt.Errorf(format,
		functionName, formatName, wordName, bitFieldName,
		bitFieldLength, bitFieldType,
	)

	return
}
