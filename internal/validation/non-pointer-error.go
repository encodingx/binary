package validation

import (
	"fmt"
)

func NewNonPointerError(functionName string) error {
	const (
		format = "" +
			"Argument to %[1]s should be a pointer to a format-struct. " +
			"Argument to %[1]s is not a pointer."
	)

	return fmt.Errorf(format, functionName)
}
