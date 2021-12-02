package errors

import (
	"fmt"
)

type formatHasNoWordsError struct {
	formatError
}

func NewFormatHasNoWordsError(formatName string) (e error) {
	e = &formatHasNoWordsError{
		formatError: formatError{
			formatName: formatName,
		},
	}

	return
}

func (e formatHasNoWordsError) Error() (message string) {
	const (
		format = "Format %s has no words in it."
	)

	message = fmt.Sprintf(format, e.formatName)

	return
}
