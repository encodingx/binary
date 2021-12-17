package validation

import (
	"fmt"
)

type FunctionError interface {
	error
	SetFunctionName(string)
}

type DefaultFunctionError struct {
	functionName string
}

func (e *DefaultFunctionError) SetFunctionName(functionName string) {
	e.functionName = functionName

	return
}

type nonPointerError struct {
	DefaultFunctionError
}

func NewNonPointerError() *nonPointerError {
	return new(nonPointerError)
}

func (e *nonPointerError) Error() string {
	const (
		format = "" +
			"Argument to %[1]s should be a pointer to a format-struct. " +
			"Argument to %[1]s is not a pointer."
	)

	return fmt.Sprintf(format, e.functionName)
}

type pointerToNonStructVariableError struct {
	DefaultFunctionError
}

func NewPointerToNonStructVariableError() *pointerToNonStructVariableError {
	return new(pointerToNonStructVariableError)
}

func (e *pointerToNonStructVariableError) Error() string {
	const (
		format = "" +
			"Argument to %[1]s should be a pointer to a format-struct. " +
			"Argument to %[1]s does not point to a struct variable."
	)

	return fmt.Sprintf(format, e.functionName)
}
