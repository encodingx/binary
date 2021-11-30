package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatIsNotStructError(t *testing.T) {
	const (
		formatName = "formatName"

		message = "Format formatName is not a struct."
	)

	var (
		e error
	)

	e = NewFormatIsNotStructError(formatName)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
