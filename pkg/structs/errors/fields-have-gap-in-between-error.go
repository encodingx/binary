package errors

import (
	"fmt"
)

type fieldsHaveGapInBetweenError struct {
	fieldsError
}

func NewFieldsHaveGapInBetweenError(
	formatName, wordName, leftFieldName, rightFieldName string,
) (
	e error,
) {
	e = &fieldsHaveGapInBetweenError{
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

func (e fieldsHaveGapInBetweenError) Error() (message string) {
	const (
		format = "There is a gap " +
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
