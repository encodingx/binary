package validation

import (
	"fmt"
)

func NewPointerToNonStructVariableError(functionName string) error {
	const (
		format = "" +
			"Argument to %[1]s should be a pointer to a format-struct. " +
			"Argument to %[1]s does not point to a struct variable."
	)

	return fmt.Errorf(format, functionName)
}
