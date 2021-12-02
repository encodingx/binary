package errors

import (
	"fmt"
)

type fieldTypeUnsupportedError struct {
	fieldError
	typeName string
}

func NewFieldTypeUnsupportedError(
	formatName, wordName, fieldName, typeName string) (
	e error,
) {
	e = &fieldTypeUnsupportedError{
		fieldError: fieldError{
			wordError: wordError{
				formatError: formatError{
					formatName: formatName,
				},
				wordName: wordName,
			},
			fieldName: fieldName,
		},
		typeName: typeName,
	}

	return
}

func (e fieldTypeUnsupportedError) Error() (message string) {
	const (
		format = "Type %s " +
			"of field %s " +
			"in word %s " +
			"in format %s " +
			"is not supported."
	)

	message = fmt.Sprintf(format,
		e.typeName,
		e.fieldName,
		e.wordName,
		e.formatName,
	)

	return
}
