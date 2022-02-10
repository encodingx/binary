package metadata

import (
	"reflect"

	"github.com/encodingx/binary/internal/validation"
)

type FormatMetadata struct {
	words         []wordMetadata
	lengthInBytes int
}

func NewFormatMetadataFromTypeReflection(reflection reflect.Type) (
	format FormatMetadata, e error,
) {
	var (
		i int
	)

	defer func() {
		if e != nil {
			e.(validation.FormatError).SetFormatName(
				reflection.String(),
			)
		}
	}()

	if reflection.NumField() == 0 {
		e = validation.NewFormatWithNoWordsError()

		return
	}

	format = FormatMetadata{
		words: make([]wordMetadata,
			reflection.NumField(),
		),
	}

	for i = 0; i < reflection.NumField(); i++ {
		format.words[i], e = newWordMetadataFromStructFieldReflection(
			reflection.Field(i),
		)
		if e != nil {
			return
		}

		format.lengthInBytes += format.words[i].lengthInBytes
	}

	return
}

func (m FormatMetadata) Marshal(reflection reflect.Value) (bytes []byte) {
	// Merge byte slices marshalled from words,
	// in the order they appear in the format.

	var (
		copyIndex int
		i         int
		word      wordMetadata
		wordBytes []byte
	)

	bytes = make([]byte, m.lengthInBytes)

	for i, word = range m.words {
		wordBytes = word.marshal(
			reflection.Field(i),
		)

		copy(bytes[copyIndex:], wordBytes)

		copyIndex += word.lengthInBytes
	}

	return
}

func (m FormatMetadata) Unmarshal(bytes []byte, reflection reflect.Value) {
	var (
		i    int
		j    int
		k    int
		word wordMetadata
	)

	for i, word = range m.words {
		k = j + word.lengthInBytes

		if k < wordLengthUpperLimitBytes {
			j = 0

		} else {
			j = k - wordLengthUpperLimitBytes
		}

		word.unmarshal(bytes[j:k],
			reflection.Field(i),
		)

		j = k
	}

	return
}

func (m FormatMetadata) LengthInBytes() int {
	return m.lengthInBytes
}
