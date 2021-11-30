package errors

import (
	"fmt"
)

type formatIsNotStructError struct {
	formatError
}

func NewFormatIsNotStructError(formatName string) (e error) {
	e = &formatIsNotStructError{
		formatError: formatError{
			formatName: formatName,
		},
	}

	return
}

func (e formatIsNotStructError) Error() (message string) {
	const (
		format = "Format %s is not a struct."
	)

	message = fmt.Sprintf(format, e.formatName)

	return
}
