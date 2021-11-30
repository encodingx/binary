package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaceIsNotPointerError(t *testing.T) {
	const (
		interfaceName = "interfaceName"

		message = "Interface interfaceName is not a pointer."
	)

	var (
		e error
	)

	e = NewInterfaceIsNotPointerError(interfaceName)

	assert.Equal(t,
		message, e.Error(),
		"The message of an error should match the expected value.",
	)
}
