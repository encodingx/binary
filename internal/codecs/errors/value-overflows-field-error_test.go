package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueOverflowsFieldError(t *testing.T) {
	const (
		fieldLengthInBits = 8
		value             = 256

		message = "Value 256 " +
			"overflows field of length 8."
	)

	var (
		e error
	)

	e = NewValueOverflowsFieldError(value, fieldLengthInBits)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
