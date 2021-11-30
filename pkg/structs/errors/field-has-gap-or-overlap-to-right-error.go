package errors

import (
	"fmt"
)

type fieldHasGapOrOverlapToRightError struct {
	fieldError
}

func NewFieldHasGapOrOverlapToRightError(
	formatName, wordName, fieldName string,
) (
	e error,
) {
	e = &fieldHasGapOrOverlapToRightError{
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

func (e fieldHasGapOrOverlapToRightError) Error() (message string) {
	const (
		format = "There is a gap or overlap " +
			"to the right of field %s " +
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
