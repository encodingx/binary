package codecs

import (
	"reflect"

	"github.com/encodingx/binary/internal/codecs/metadata"
	"github.com/encodingx/binary/internal/validation"
)

type Codec struct {
	formatMetadataCache map[reflect.Type]metadata.FormatMetadata
}

func NewCodec() (c Codec) {
	c = Codec{
		formatMetadataCache: make(map[reflect.Type]metadata.FormatMetadata),
	}

	return
}

func (c Codec) formatMetadataFromTypeReflection(reflection reflect.Type) (
	format metadata.FormatMetadata, e error,
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

	format, e = metadata.NewFormatMetadataFromTypeReflection(
		reflection.Elem(),
	)
	if e != nil {
		return
	}

	c.formatMetadataCache[reflection] = format

	return
}

func (c Codec) NewOperation(iface interface{}) (
	operation CodecOperation, e error,
) {
	operation.format, e = c.formatMetadataFromTypeReflection(
		reflect.TypeOf(iface),
	)
	if e != nil {
		return
	}

	operation.valueReflection = reflect.ValueOf(iface).Elem()

	return
}

type CodecOperation struct {
	format          metadata.FormatMetadata
	valueReflection reflect.Value
}

func (c CodecOperation) Marshal() (bytes []byte, e error) {
	bytes = c.format.Marshal(c.valueReflection)

	return
}

func (c CodecOperation) Unmarshal(bytes []byte) (e error) {
	if len(bytes) != c.format.LengthInBytes() {
		e = validation.NewLengthOfByteSliceNotEqualToFormatLengthError(
			uint(c.format.LengthInBytes()),
			uint(len(bytes)),
		)

		e.(validation.FormatError).SetFormatName(
			c.valueReflection.Type().String(),
		)

		return
	}

	c.format.Unmarshal(bytes, c.valueReflection)

	return
}
