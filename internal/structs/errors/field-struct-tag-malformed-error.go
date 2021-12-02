package errors

import (
	"fmt"
)

type fieldStructTagMalformedError struct {
	fieldError
}

func NewFieldStructTagMalformedError(formatName, wordName, fieldName string) (
	e error,
) {
	e = &fieldStructTagMalformedError{
		fieldError: fieldError{
			wordError: wordError{
				formatError: formatError{
					formatName: formatName,
				},
				wordName: wordName,
			},
			fieldName: fieldName,
		},
	}

	return
}

func (e fieldStructTagMalformedError) Error() (message string) {
	const (
		format = "Field %s " +
			"in word %s " +
			"in format %s " +
			"has a malformed struct tag."
	)

	message = fmt.Sprintf(format,
		e.fieldName,
		e.wordName,
		e.formatName,
	)

	return
}
