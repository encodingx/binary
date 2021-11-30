package codecs

import (
	"encoding/binary"
	"reflect"

	"github.com/joel-ling/go-bitfields/pkg/codecs/errors"
	"github.com/joel-ling/go-bitfields/pkg/structs"
	"github.com/joel-ling/go-bitfields/pkg/structs/formats"
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
		bytesI []byte
		format formats.Format
		i      uint
		j      uint
		l      uint
		length uint
		offset uint
		valueI uint64
		valueJ uint64
		values reflect.Value
	)

	format, e = c.formatStructParser.ParseFormatStruct(pointer)
	if e != nil {
		return
	}

	bytes = make([]byte,
		format.LengthInBytes(),
	)

	values = reflect.ValueOf(pointer).Elem()

	for i = 0; i < format.NWords(); i++ {
		valueI = 0

		for j = 0; j < format.Word(i).NFields(); j++ {
			switch format.Word(i).Field(j).Kind() {
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fallthrough

			case reflect.Uint:
				valueJ = values.Field(int(i)).Field(int(j)).Uint()

			case reflect.Bool:
				if values.Field(int(i)).Field(int(j)).Bool() {
					valueJ = 1
				}
			}

			length = format.Word(i).Field(j).LengthInBits()
			offset = format.Word(i).Field(j).OffsetInBits()

			if valueJ > (1<<length - 1) {
				e = errors.NewValueOverflowsFieldError(uint(valueJ), length)

				return
			}

			valueI = (valueJ << offset) | valueI
		}

		bytesI = make([]byte, maximumWordLengthInBytes)

		binary.BigEndian.PutUint64(bytesI, valueI)

		copy(bytes[l:],
			bytesI[maximumWordLengthInBytes-format.Word(i).LengthInBytes():],
		)

		l += format.Word(i).LengthInBytes()
	}

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
