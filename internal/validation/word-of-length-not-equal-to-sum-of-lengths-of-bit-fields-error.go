package validation

import (
	"fmt"
)

func NewWordOfLengthNotEqualToSumOfLengthsOfBitFieldsError(
	functionName, formatName, wordName string,
	wordLength, bitFieldLengthSum uint,
) (
	e error,
) {
	const (
		format = "" +
			"The length of a word " +
			"should be equal to the sum of lengths of its bit fields. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"that has a word \"%s\" " +
			"of length %d " +
			"not equal to the sum of the lengths of its bit fields, %d."
	)

	e = fmt.Errorf(format,
		functionName, formatName, wordName,
		wordLength, bitFieldLengthSum,
	)

	return
}
