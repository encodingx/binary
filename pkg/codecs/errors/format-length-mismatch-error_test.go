package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatLengthMismatchError(t *testing.T) {
	const (
		actualLengthInBytes   = 9
		expectedLengthInBytes = 8

		message = "Length 9 of byte slice " +
			"does not match format length of 8 bytes."
	)

	var (
		e error
	)

	e = NewFormatLengthMismatchError(expectedLengthInBytes, actualLengthInBytes)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
