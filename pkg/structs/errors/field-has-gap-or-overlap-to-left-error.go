package errors

import (
	"fmt"
)

type fieldHasGapOrOverlapToLeftError struct {
	fieldError
}

func NewFieldHasGapOrOverlapToLeftError(
	formatName, wordName, fieldName string,
) (
	e error,
) {
	e = &fieldHasGapOrOverlapToLeftError{
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

func (e fieldHasGapOrOverlapToLeftError) Error() (message string) {
	const (
		format = "There is a gap or overlap " +
			"to the left of field %s " +
			"in word %s " +
			"in format %s."
	)

	message = fmt.Sprintf(format,
		e.fieldName,
		e.wordName,
		e.formatName,
	)

	return
}
