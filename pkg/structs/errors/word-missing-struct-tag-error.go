package errors

import (
	"fmt"
)

type wordMissingStructTagError struct {
	wordError
}

func NewWordMissingStructTagError(formatName, wordName string) (e error) {
	e = &wordMissingStructTagError{
		wordError: wordError{
			formatName: formatName,
			wordName:   wordName,
		},
	}

	return
}

func (e wordMissingStructTagError) Error() (message string) {
	const (
		format = "Word %s " +
			"in format %s " +
			"is missing a struct tag."
	)

	message = fmt.Sprintf(format,
		e.wordName,
		e.formatName,
	)

	return
}
