package errors

import (
	"fmt"
)

type wordIsNotStructError struct {
	wordError
}

func NewWordIsNotStructError(formatName, wordName string) (e error) {
	e = &wordIsNotStructError{
		wordError: wordError{
			formatError: formatError{
				formatName: formatName,
			},
			wordName: wordName,
		},
	}

	return
}

func (e wordIsNotStructError) Error() (message string) {
	const (
		format = "Word %s " +
			"in format %s " +
			"is not a struct."
	)

	message = fmt.Sprintf(format,
		e.wordName,
		e.formatName,
	)

	return
}
