package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	functionName = "Marshal"
)

func TestNonPointerError(t *testing.T) {
	const (
		errorMessage = "" +
			"Argument to Marshal should be a pointer to a format-struct. " +
			"Argument to Marshal is not a pointer."
	)

	var (
		e FunctionError
	)

	e = NewNonPointerError()

	e.SetFunctionName(functionName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}

func TestPointerToNonStructVariableError(t *testing.T) {
	const (
		errorMessage = "" +
			"Argument to Marshal should be a pointer to a format-struct. " +
			"Argument to Marshal does not point to a struct variable."
	)

	var (
		e FunctionError
	)

	e = NewPointerToNonStructVariableError()

	e.SetFunctionName(functionName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
