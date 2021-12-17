package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointerToNonStructVariableError(t *testing.T) {
	const (
		functionName = "Marshal"

		errorMessage = "" +
			"Argument to Marshal should be a pointer to a format-struct. " +
			"Argument to Marshal does not point to a struct variable."
	)

	var (
		e error
	)

	e = NewPointerToNonStructVariableError(functionName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
