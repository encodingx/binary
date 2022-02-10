package metadata

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/encodingx/binary/internal/validation"
)

type wordMetadata struct {
	bitFields     []bitFieldMetadata
	lengthInBits  uint
	lengthInBytes int
}

func newWordMetadataFromStructFieldReflection(reflection reflect.StructField) (
	word wordMetadata, e error,
) {
	const (
		tagKey         = "word"
		tagValueFormat = "%d"

		wordLengthFactor     = 8
		wordLengthLowerLimit = 8
		wordLengthUpperLimit = 64
	)

	var (
		offset       uint
		wordLength   uint
		wordLengthOK bool

		i int
	)

	defer func() {
		if e != nil {
			e.(validation.WordError).SetWordName(reflection.Name)
		}
	}()

	if reflection.Type.Kind() != reflect.Struct {
		e = validation.NewWordNotStructError()

		return
	}

	if len(reflection.Tag) == 0 {
		e = validation.NewWordWithNoStructTagError()

		return
	}

	_, e = fmt.Sscanf(
		reflection.Tag.Get(tagKey),
		tagValueFormat,
		&wordLength,
	)
	if e != nil {
		e = validation.NewWordWithMalformedTagError()

		return
	}

	wordLengthOK = wordLength%wordLengthFactor == 0
	wordLengthOK = wordLengthOK && wordLength >= wordLengthLowerLimit
	wordLengthOK = wordLengthOK && wordLength <= wordLengthUpperLimit

	if !wordLengthOK {
		e = validation.NewWordOfIncompatibleLengthError(wordLength)

		return
	}

	if reflection.Type.NumField() == 0 {
		e = validation.NewWordWithNoBitFieldsError()

		return
	}

	word = wordMetadata{
		bitFields: make([]bitFieldMetadata,
			reflection.Type.NumField(),
		),
		lengthInBits:  wordLength,
		lengthInBytes: int(wordLength / wordLengthFactor),
	}

	offset = wordLength

	for i = 0; i < reflection.Type.NumField(); i++ {
		word.bitFields[i], e = newBitFieldMetadataFromStructFieldReflection(
			reflection.Type.Field(i),
		)
		if e != nil {
			return
		}

		offset -= word.bitFields[i].length

		word.bitFields[i].offset = uint64(offset)
	}

	if offset != 0 {
		e = validation.NewWordOfLengthNotEqualToSumOfLengthsOfBitFieldsError(
			wordLength,
			wordLength-offset,
		)

		return
	}

	return
}

func (m wordMetadata) marshal(reflection reflect.Value) (bytes []byte) {
	var (
		bitField       bitFieldMetadata
		bitFieldUint64 uint64
		i              int
		wordUint64     uint64
	)

	for i, bitField = range m.bitFields {
		bitFieldUint64 = bitField.marshal(
			reflection.Field(i),
		)

		wordUint64 = wordUint64 | bitFieldUint64
	}

	bytes = make([]byte, wordLengthUpperLimitBytes)

	binary.BigEndian.PutUint64(bytes, wordUint64)

	bytes = bytes[wordLengthUpperLimitBytes-m.lengthInBytes:]

	return
}

func (m wordMetadata) unmarshal(bytes []byte, reflection reflect.Value) {
	var (
		bitFieldBytes []byte
		i             int
	)

	for i = 0; i < len(m.bitFields); i++ {
		if len(bytes) < wordLengthUpperLimitBytes {
			bitFieldBytes = make([]byte, wordLengthUpperLimitBytes)

			copy(bitFieldBytes[wordLengthUpperLimitBytes-len(bytes):],
				bytes,
			)

		} else {
			bitFieldBytes = bytes
		}

		m.bitFields[i].unmarshal(bitFieldBytes,
			reflection.Field(i),
		)
	}

	return
}
