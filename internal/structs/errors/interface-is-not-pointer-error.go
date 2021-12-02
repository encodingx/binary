package errors

import (
	"fmt"
)

type interfaceIsNotPointerError struct {
	interfaceName string
}

func NewInterfaceIsNotPointerError(interfaceName string) (e error) {
	e = &interfaceIsNotPointerError{
		interfaceName: interfaceName,
	}

	return
}

func (e interfaceIsNotPointerError) Error() (message string) {
	const (
		format = "Interface %s is not a pointer."
	)

	message = fmt.Sprintf(format, e.interfaceName)

	return
}
