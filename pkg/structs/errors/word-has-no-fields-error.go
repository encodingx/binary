package errors

import (
	"fmt"
)

type wordHasNoFieldsError struct {
	wordError
}

func NewWordHasNoFieldsError(formatName, wordName string) (e error) {
	e = &wordHasNoFieldsError{
		wordError: wordError{
			formatError: formatError{
				formatName: formatName,
			},
			wordName: wordName,
		},
	}

	return
}

func (e wordHasNoFieldsError) Error() (message string) {
	const (
		format = "Word %s " +
			"in format %s " +
			"has no fields in it."
	)

	message = fmt.Sprintf(format,
		e.wordName,
		e.formatName,
	)

	return
}
