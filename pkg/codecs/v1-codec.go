package codecs

import (
	"encoding/binary"
	"reflect"

	"github.com/joel-ling/go-bitfields/pkg/codecs/errors"
	"github.com/joel-ling/go-bitfields/pkg/structs"
	"github.com/joel-ling/go-bitfields/pkg/structs/fields"
	"github.com/joel-ling/go-bitfields/pkg/structs/formats"
	"github.com/joel-ling/go-bitfields/pkg/structs/words"
)

const (
	maximumWordLengthInBytes = 8
)

type v1Codec struct {
	formatStructParser structs.FormatStructParser
}

func NewV1Codec() (codec *v1Codec) {
	// Return a Version 1 implementation of interface Codec.

	codec = &v1Codec{
		formatStructParser: structs.NewFormatStructParser(),
	}

	return
}

func (c *v1Codec) Marshal(pointer interface{}) (bytes []byte, e error) {
	// Implement interface Codec.

	var (
		format formats.Format
	)

	format, e = c.formatStructParser.ParseFormatStruct(pointer)
	if e != nil {
		return
	}

	bytes, e = marshalFormat(format,
		reflect.ValueOf(pointer).Elem(),
	)

	return
}

func marshalFormat(format formats.Format, valueReflection reflect.Value) (
	bytes []byte, e error,
) {
	// Merge byte slices containing words,
	// in the order they appear in the format.

	var (
		bytesI []byte
		i      uint
		j      uint
	)

	bytes = make([]byte,
		format.LengthInBytes(),
	)

	for i = 0; i < format.NWords(); i++ {
		bytesI, e = marshalWord(
			format.Word(i),
			valueReflection.Field(int(i)),
		)
		if e != nil {
			return
		}

		copy(bytes[j:], bytesI)

		j += format.Word(i).LengthInBytes()
	}

	return
}

func marshalWord(word words.Word, valueReflection reflect.Value) (
	bytes []byte, e error,
) {
	// Accumulate the bitwise OR sum of the bit-shifted values
	// of fields in the word,
	// convert that to a right-aligned sequence of bytes and
	// trim off the leftmost bytes in excess of word length.

	var (
		i      uint
		value  uint64
		valueI uint64
	)

	for i = 0; i < word.NFields(); i++ {
		valueI, e = marshalField(
			word.Field(i),
			valueReflection.Field(int(i)),
		)
		if e != nil {
			return
		}

		value = value | valueI
	}

	bytes = make([]byte, maximumWordLengthInBytes)

	binary.BigEndian.PutUint64(bytes, value)

	bytes = bytes[maximumWordLengthInBytes-word.LengthInBytes():]

	return
}

func marshalField(field fields.Field, valueReflection reflect.Value) (
	value uint64, e error,
) {
	// Derive value of field through reflection,
	// convert booleans to integers,
	// ensure that value does not overflow field, and
	// return the value bit-shifted.

	switch field.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough

	case reflect.Uint:
		value = valueReflection.Uint()

	case reflect.Bool:
		if valueReflection.Bool() {
			value = 1
		}
	}

	if value > (1<<field.LengthInBits() - 1) {
		e = errors.NewValueOverflowsFieldError(
			uint(value),
			field.LengthInBits(),
		)

		return
	}

	value = value << field.OffsetInBits()

	return
}

func (c *v1Codec) Unmarshal(bytes []byte, pointer interface{}) (e error) {
	// Implement interface Codec.

	var (
		bytesI []byte
		bytesJ []byte
		format formats.Format
		i      uint
		j      uint
		k      uint
		l      uint
		length uint
		offset uint
		value  uint64
		values reflect.Value
	)

	format, e = c.formatStructParser.ParseFormatStruct(pointer)
	if e != nil {
		return
	}

	if uint(len(bytes)) != format.LengthInBytes() {
		e = errors.NewFormatLengthMismatchError(
			format.LengthInBytes(),
			uint(len(bytes)),
		)

		return
	}

	values = reflect.ValueOf(pointer).Elem()

	for i = 0; i < format.NWords(); i++ {
		l = k + format.Word(i).LengthInBytes()

		bytesI = bytes[k:l]

		k = l

		for j = 0; j < format.Word(i).NFields(); j++ {
			bytesJ = make([]byte, maximumWordLengthInBytes)

			copy(bytesJ[maximumWordLengthInBytes-format.Word(i).LengthInBytes():],
				bytesI,
			)

			length = format.Word(i).Field(j).LengthInBits()
			offset = format.Word(i).Field(j).OffsetInBits()

			value = binary.BigEndian.Uint64(bytesJ) >> offset & (1<<length - 1)

			switch format.Word(i).Field(j).Kind() {
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fallthrough

			case reflect.Uint:
				values.Field(int(i)).Field(int(j)).SetUint(value)

			case reflect.Bool:
				if value == 1 {
					values.Field(int(i)).Field(int(j)).SetBool(true)
				}
			}
		}

	}

	return
}
