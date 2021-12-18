package binary

import (
	"reflect"

	"github.com/encodingx/binary/internal/validation"
)

type codec struct {
	formatMetadataCache map[reflect.Type]*formatMetadata
}

func newCodec() (c *codec) {
	c = &codec{
		formatMetadataCache: make(map[reflect.Type]*formatMetadata),
	}

	return
}

func (c *codec) formatMetadataFromTypeReflection(reflection reflect.Type) (
	format *formatMetadata, e error,
) {
	var (
		inCache bool
	)

	format, inCache = c.formatMetadataCache[reflection]

	if inCache {
		return
	}

	if reflection.Kind() != reflect.Ptr {
		e = validation.NewNonPointerError()

		return
	}

	if reflection.Elem().Kind() != reflect.Struct {
		e = validation.NewPointerToNonStructVariableError()

		return
	}

	format, e = newFormatMetadataFromTypeReflection(
		reflection.Elem(),
	)
	if e != nil {
		return
	}

	c.formatMetadataCache[reflection] = format

	return
}

func (c *codec) newOperation(iface interface{}) (
	operation *codecOperation, e error,
) {
	operation = new(codecOperation)

	operation.format, e = c.formatMetadataFromTypeReflection(
		reflect.TypeOf(iface),
	)
	if e != nil {
		return
	}

	operation.valueReflection = reflect.ValueOf(iface).Elem()

	return
}

type codecOperation struct {
	format          *formatMetadata
	valueReflection reflect.Value
}

func (c *codecOperation) marshal() (bytes []byte, e error) {
	bytes = c.format.marshal(c.valueReflection)

	return
}

func (c *codecOperation) unmarshal(bytes []byte) (e error) {
	if len(bytes) != c.format.lengthInBytes {
		e = validation.NewLengthOfByteSliceNotEqualToFormatLengthError(
			uint(c.format.lengthInBytes),
			uint(len(bytes)),
		)

		e.(validation.FormatError).SetFormatName(
			c.valueReflection.Type().String(),
		)

		return
	}

	c.format.unmarshal(bytes, c.valueReflection)

	return
}
