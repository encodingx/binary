package binary

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/encodingx/binary/internal/validation"
)

func Marshal(iface interface{}) (bytes []byte, e error) {
	const (
		functionName = "Marshal"
	)

	var (
		format          *formatMetadata
		typeReflection  reflect.Type
		valueReflection reflect.Value
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

	typeReflection, valueReflection, e = structReflectionFromInterface(iface)
	if e != nil {
		return
	}

	format, e = newFormatMetadataFromTypeReflection(typeReflection)
	if e != nil {
		return
	}

	bytes = format.marshal(valueReflection)

	return
}

func Unmarshal(bytes []byte, iface interface{}) (e error) {
	const (
		functionName = "Unmarshal"
	)

	var (
		format     *formatMetadata
		reflection reflect.Type
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

	reflection, _, e = structReflectionFromInterface(iface)
	if e != nil {
		return
	}

	format, e = newFormatMetadataFromTypeReflection(reflection)
	if e != nil {
		return
	}

	if len(bytes) != format.lengthInBytes {
		e = validation.NewLengthOfByteSliceNotEqualToFormatLengthError(
			uint(format.lengthInBytes),
			uint(len(bytes)),
		)

		e.(validation.FormatError).SetFormatName(
			reflection.String(),
		)

		return
	}

	return
}

func structReflectionFromInterface(iface interface{}) (
	typeReflection reflect.Type, valueReflection reflect.Value, e error,
) {
	typeReflection = reflect.TypeOf(iface)

	if typeReflection.Kind() != reflect.Ptr {
		e = validation.NewNonPointerError()

		return
	}

	typeReflection = typeReflection.Elem()

	if typeReflection.Kind() != reflect.Struct {
		e = validation.NewPointerToNonStructVariableError()

		return
	}

	valueReflection = reflect.ValueOf(iface).Elem()

	return
}

type formatMetadata struct {
	words         []*wordMetadata
	lengthInBytes int
}

func newFormatMetadataFromTypeReflection(reflection reflect.Type) (
	format *formatMetadata, e error,
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

	format = &formatMetadata{
		words: make([]*wordMetadata,
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

func (m *formatMetadata) marshal(reflection reflect.Value) (bytes []byte) {
	var (
		copyIndex int
		i         int
		word      *wordMetadata
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

type wordMetadata struct {
	bitfields     []*bitFieldMetadata
	lengthInBits  uint
	lengthInBytes int
}

func newWordMetadataFromStructFieldReflection(reflection reflect.StructField) (
	word *wordMetadata, e error,
) {
	const (
		tagKey         = "word"
		tagValueFormat = "%d"

		wordLengthFactor     = 8
		wordLengthLowerLimit = 8
		wordLengthUpperLimit = 64
	)

	var (
		sumBitFieldLengths uint
		wordLength         uint
		wordLengthOK       bool

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

	word = &wordMetadata{
		bitfields: make([]*bitFieldMetadata,
			reflection.Type.NumField(),
		),
		lengthInBits:  wordLength,
		lengthInBytes: int(wordLength / wordLengthFactor),
	}

	for i = 0; i < reflection.Type.NumField(); i++ {
		word.bitfields[i], e = newBitFieldMetadataFromStructFieldReflection(
			reflection.Type.Field(i),
		)
		if e != nil {
			return
		}

		sumBitFieldLengths += word.bitfields[i].length
	}

	if sumBitFieldLengths != wordLength {
		e = validation.NewWordOfLengthNotEqualToSumOfLengthsOfBitFieldsError(
			wordLength,
			sumBitFieldLengths,
		)

		return
	}

	return
}

func (m *wordMetadata) marshal(reflection reflect.Value) (bytes []byte) {
	const (
		wordLengthUpperLimitBytes = 8
	)

	var (
		bitField       *bitFieldMetadata
		bitFieldUint64 uint64
		i              int
		offset         uint
		wordUint64     uint64
	)

	offset = m.lengthInBits

	for i, bitField = range m.bitfields {
		offset -= bitField.length

		bitFieldUint64 = bitField.marshal(
			reflection.Field(i),
			uint64(offset),
		)

		wordUint64 = wordUint64 | bitFieldUint64
	}

	bytes = make([]byte, wordLengthUpperLimitBytes)

	binary.BigEndian.PutUint64(bytes, wordUint64)

	bytes = bytes[wordLengthUpperLimitBytes-m.lengthInBytes:]

	return
}

type bitFieldMetadata struct {
	length uint
	kind   reflect.Kind
}

func newBitFieldMetadataFromStructFieldReflection(
	reflection reflect.StructField,
) (
	bitField *bitFieldMetadata, e error,
) {
	const (
		tagKey         = "bitfield"
		tagValueFormat = "%d"
	)

	var (
		bitFieldLengthCap uint
	)

	defer func() {
		if e != nil {
			e.(validation.BitFieldError).SetBitFieldName(reflection.Name)
		}
	}()

	switch reflection.Type.Kind() {
	case reflect.Uint:
		fallthrough

	case reflect.Uint64:
		bitFieldLengthCap = 64

	case reflect.Uint32:
		bitFieldLengthCap = 32

	case reflect.Uint16:
		bitFieldLengthCap = 16

	case reflect.Uint8:
		bitFieldLengthCap = 8

	case reflect.Bool:
		bitFieldLengthCap = 1

	default:
		e = validation.NewBitFieldOfUnsupportedTypeError(
			reflection.Type.String(),
		)

		return
	}

	bitField = &bitFieldMetadata{
		kind: reflection.Type.Kind(),
	}

	if len(reflection.Tag) == 0 {
		e = validation.NewBitFieldWithNoStructTagError()

		return
	}

	_, e = fmt.Sscanf(
		reflection.Tag.Get(tagKey),
		tagValueFormat,
		&bitField.length,
	)
	if e != nil {
		e = validation.NewBitFieldWithMalformedTagError()

		return
	}

	if bitField.length > bitFieldLengthCap {
		e = validation.NewBitFieldOfLengthOverflowingTypeError(
			bitField.length,
			reflection.Type.String(),
		)
	}

	return
}

func (m *bitFieldMetadata) marshal(
	reflection reflect.Value, offset uint64,
) (
	value uint64,
) {
	switch m.kind {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough

	case reflect.Uint:
		value = reflection.Uint()

	case reflect.Bool:
		if reflection.Bool() {
			value = 1
		}
	}

	value = value & (1<<m.length - 1) // XXX: mask if overflowing

	value = value << offset

	return
}
