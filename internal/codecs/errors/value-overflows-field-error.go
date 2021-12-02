package errors

import (
	"fmt"
)

type valueOverflowsFieldError struct {
	value             uint
	fieldLengthInBits uint
}

func NewValueOverflowsFieldError(value, fieldLengthInBits uint) (e error) {
	e = &valueOverflowsFieldError{
		value:             value,
		fieldLengthInBits: fieldLengthInBits,
	}

	return
}

func (e valueOverflowsFieldError) Error() (message string) {
	const (
		format = "Value %d " +
			"overflows field of length %d."
	)

	message = fmt.Sprintf(format,
		e.value,
		e.fieldLengthInBits,
	)

	return
}
