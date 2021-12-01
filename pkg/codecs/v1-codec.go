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
		format formats.Format
	)

	format, e = c.formatStructParser.ParseFormatStruct(pointer)
	if e != nil {
		return
	}

	e = unmarshalFormat(format,
		reflect.ValueOf(pointer).Elem(),
		bytes,
	)
	if e != nil {
		return
	}

	return
}

func unmarshalFormat(
	format formats.Format, valueReflection reflect.Value, bytes []byte,
) (
	e error,
) {
	// Ensure that the length of the byte slice matches the format, and
	// dispatch sections of the slice to be unmarshalled into words.

	var (
		i uint
		j uint
		k uint
	)

	if uint(len(bytes)) != format.LengthInBytes() {
		e = errors.NewFormatLengthMismatchError(
			format.LengthInBytes(),
			uint(len(bytes)),
		)

		return
	}

	for i = 0; i < format.NWords(); i++ {
		k = j + format.Word(i).LengthInBytes()

		unmarshalWord(
			format.Word(i),
			valueReflection.Field(int(i)),
			bytes[j:k],
		)

		j = k
	}

	return
}

func unmarshalWord(
	word words.Word, valueReflection reflect.Value, bytes []byte,
) {
	// Pad the bytes of a word into a right-aligned slice of fixed length, and
	// unmarshal every field in the word from that slice.

	var (
		bytesI []byte
		i      uint
	)

	for i = 0; i < word.NFields(); i++ {
		bytesI = make([]byte, maximumWordLengthInBytes)

		copy(
			bytesI[maximumWordLengthInBytes-word.LengthInBytes():],
			bytes,
		)

		unmarshalField(
			word.Field(i),
			valueReflection.Field(int(i)),
			bytesI,
		)
	}

	return
}

func unmarshalField(
	field fields.Field, valueReflection reflect.Value, bytes []byte,
) {
	// Recover the value of a field from a bitwise OR sum of bit-shifted values
	// by reversing the bit shift and
	// performing a bitwise AND with the bit mask of the field.
	// Convert integers to booleans if applicable and set values by reflection.

	var (
		value uint64
	)

	value = binary.BigEndian.Uint64(bytes) >>
		field.OffsetInBits() & (1<<field.LengthInBits() - 1)

	switch field.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough

	case reflect.Uint:
		valueReflection.SetUint(value)

	case reflect.Bool:
		switch value {
		case 1:
			valueReflection.SetBool(true)
		default:
			valueReflection.SetBool(false)
		}
	}

	return
}
