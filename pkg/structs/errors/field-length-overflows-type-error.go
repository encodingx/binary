package errors

import (
	"fmt"
)

type fieldLengthOverflowsTypeError struct {
	fieldError
	typeName    string
	typeLength  uint
	fieldLength uint
}

func NewFieldLengthOverflowsTypeError(
	formatName, wordName, fieldName, typeName string,
	typeLength, fieldLength uint,
) (
	e error,
) {
	e = &fieldLengthOverflowsTypeError{
		fieldError: fieldError{
			wordError: wordError{
				formatError: formatError{
					formatName: formatName,
				},
				wordName: wordName,
			},
			fieldName: fieldName,
		},
		typeName:    typeName,
		typeLength:  typeLength,
		fieldLength: fieldLength,
	}

	return
}

func (e fieldLengthOverflowsTypeError) Error() (message string) {
	const (
		format = "Length %d bits " +
			"of field %s " +
			"in word %s " +
			"in format %s " +
			"overflows type %s " +
			"of length %d bits."
	)

	message = fmt.Sprintf(format,
		e.fieldLength,
		e.fieldName,
		e.wordName,
		e.formatName,
		e.typeName,
		e.typeLength,
	)

	return
}
