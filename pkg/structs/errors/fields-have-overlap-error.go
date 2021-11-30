package errors

import (
	"fmt"
)

type fieldsHaveOverlapError struct {
	fieldsError
}

func NewFieldsHaveOverlapError(
	formatName, wordName, leftFieldName, rightFieldName string,
) (
	e error,
) {
	e = &fieldsHaveOverlapError{
		fieldsError: fieldsError{
			wordError: wordError{
				formatError: formatError{
					formatName: formatName,
				},
				wordName: wordName,
			},
			leftFieldName:  leftFieldName,
			rightFieldName: rightFieldName,
		},
	}

	return
}

func (e fieldsHaveOverlapError) Error() (message string) {
	const (
		format = "There is an overlap " +
			"between fields %s " +
			"and %s " +
			"in word %s " +
			"in format %s."
	)

	message = fmt.Sprintf(format,
		e.leftFieldName,
		e.rightFieldName,
		e.wordName,
		e.formatName,
	)

	return
}
