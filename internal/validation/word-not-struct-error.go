package validation

import (
	"fmt"
)

func NewWordNotStructError(functionName, formatName, wordName string) error {
	const (
		format = "" +
			"A format-struct should nest exported word-structs. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"that is not a struct."
	)

	return fmt.Errorf(format, functionName, formatName, wordName)
}
