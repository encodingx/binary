package validation

import (
	"fmt"
)

func NewWordWithNoBitFieldsError(
	functionName, formatName, wordName string,
) (
	e error,
) {
	const (
		format = "" +
			"A word-struct should have exported fields " +
			"representing bit fields. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"that has no bit fields."
	)

	return fmt.Errorf(format, functionName, formatName, wordName)
}
