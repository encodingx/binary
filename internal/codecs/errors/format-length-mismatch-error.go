package errors

import (
	"fmt"
)

type formatLengthMismatchError struct {
	expectedLengthInBytes uint
	actualLengthInBytes   uint
}

func NewFormatLengthMismatchError(
	expectedLengthInBytes, actualLengthInBytes uint,
) (
	e error,
) {
	e = &formatLengthMismatchError{
		expectedLengthInBytes: expectedLengthInBytes,
		actualLengthInBytes:   actualLengthInBytes,
	}

	return
}

func (e formatLengthMismatchError) Error() (message string) {
	const (
		format = "Length %d of byte slice " +
			"does not match format length of %d bytes."
	)

	message = fmt.Sprintf(format,
		e.actualLengthInBytes,
		e.expectedLengthInBytes,
	)

	return
}
