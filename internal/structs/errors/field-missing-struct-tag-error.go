package errors

import (
	"fmt"
)

type fieldMissingStructTagError struct {
	fieldError
}

func NewFieldMissingStructTagError(formatName, wordName, fieldName string) (
	e error,
) {
	e = &fieldMissingStructTagError{
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

func (e fieldMissingStructTagError) Error() (message string) {
	const (
		format = "Field %s " +
			"in word %s " +
			"in format %s " +
			"is missing a struct tag."
	)

	message = fmt.Sprintf(format,
		e.fieldName,
		e.wordName,
		e.formatName,
	)

	return
}
