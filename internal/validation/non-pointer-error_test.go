package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonPointerError(t *testing.T) {
	const (
		functionName = "Marshal"

		errorMessage = "" +
			"Argument to Marshal should be a pointer to a format-struct. " +
			"Argument to Marshal is not a pointer."
	)

	var (
		e error
	)

	e = NewNonPointerError(functionName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
