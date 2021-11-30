package errors

import (
	"fmt"
)

type wordStructTagMalformedError struct {
	wordError
}

func NewWordStructTagMalformedError(formatName, wordName string) (e error) {
	e = &wordStructTagMalformedError{
		wordError: wordError{
			formatName: formatName,
			wordName:   wordName,
		},
	}

	return
}

func (e wordStructTagMalformedError) Error() (message string) {
	const (
		format = "Word %s " +
			"in format %s " +
			"has a malformed struct tag."
	)

	message = fmt.Sprintf(format,
		e.wordName,
		e.formatName,
	)

	return
}
