package validation

import (
	"fmt"
)

func NewWordOfIncompatibleLengthError(
	functionName, formatName, wordName string, wordLength uint,
) (
	e error,
) {
	const (
		format = "" +
			"The length of a word should be a multiple of eight " +
			"in the range [8, 64]. " +
			"Argument to %[1]s points to a format-struct \"%s\" " +
			"that has a word \"%s\" " +
			"of length %d not in {8, 16, 24, ... 64}."
	)

	return fmt.Errorf(format, functionName, formatName, wordName, wordLength)
}
