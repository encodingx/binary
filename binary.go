package binary

import (
	"fmt"

	"github.com/encodingx/binary/internal/validation"
)

const (
	wordLengthUpperLimitBytes = 8
)

var (
	defaultCodec = newCodec()
)

func Marshal(iface interface{}) (bytes []byte, e error) {
	const (
		functionName = "Marshal"
	)

	var (
		operation *codecOperation
	)

	defer func() {
		const (
			marshalError = "Marshal error: %w"
		)

		if e != nil {
			e.(validation.FunctionError).SetFunctionName(functionName)

			e = fmt.Errorf(marshalError, e)
		}

		return
	}()

	operation, e = defaultCodec.newOperation(iface)
	if e != nil {
		return
	}

	bytes, e = operation.marshal()
	if e != nil {
		return
	}

	return
}

func Unmarshal(bytes []byte, iface interface{}) (e error) {
	const (
		functionName = "Unmarshal"
	)

	var (
		operation *codecOperation
	)

	defer func() {
		const (
			unmarshalError = "Unmarshal error: %w"
		)

		if e != nil {
			e.(validation.FunctionError).SetFunctionName(functionName)

			e = fmt.Errorf(unmarshalError, e)
		}

		return
	}()

	operation, e = defaultCodec.newOperation(iface)
	if e != nil {
		return
	}

	e = operation.unmarshal(bytes)
	if e != nil {
		return
	}

	return
}
