package validation

import (
	"fmt"
)

func NewFormatWithNoWordsError(functionName, formatName string) error {
	const (
		format = "" +
			"A format-struct should nest exported word-structs. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"that has no words."
	)

	return fmt.Errorf(format, functionName, formatName)
}
