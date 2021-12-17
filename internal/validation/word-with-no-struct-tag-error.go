package validation

import (
	"fmt"
)

func NewWordWithNoStructTagError(
	functionName, formatName, wordName string,
) (
	e error,
) {
	const (
		format = "" +
			"A format-struct should nest exported word-structs " +
			"tagged with a key \"word\" and a value " +
			"indicating the length of a word in number of bits " +
			"(e.g. `word:\"32\"`). " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"with no struct tag."
	)

	return fmt.Errorf(format, functionName, formatName, wordName)
}
