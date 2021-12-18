package binary

import (
	"encoding/binary"
	"fmt"
	"reflect"

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

	format, e = newFormatMetadataFromTypeReflection(reflection)
	if e != nil {
		return
	}

	c.formatMetadataCache[reflection] = format

	return
}

func (c *codec) newOperation(iface interface{}) (
	operation *codecOperation, e error,
) {
	var (
		reflection reflect.Type
	)

	reflection = reflect.TypeOf(iface)

	if reflection.Kind() != reflect.Ptr {
		e = validation.NewNonPointerError()

		return
	}

	reflection = reflection.Elem()

	if reflection.Kind() != reflect.Struct {
		e = validation.NewPointerToNonStructVariableError()

		return
	}

	operation = new(codecOperation)

	operation.format, e = c.formatMetadataFromTypeReflection(reflection)
	if e != nil {
		return
	}

	operation.formatName = reflection.String()

	operation.valueReflection = reflect.ValueOf(iface).Elem()

	return
}

type codecOperation struct {
	format          *formatMetadata
	formatName      string
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

		e.(validation.FormatError).SetFormatName(c.formatName)

		return
	}

	c.format.unmarshal(bytes, c.valueReflection)

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

func (m *formatMetadata) unmarshal(bytes []byte, reflection reflect.Value) {
	var (
		i    int
		j    int
		k    int
		word *wordMetadata
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

type wordMetadata struct {
	bitFields     []*bitFieldMetadata
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

	word = &wordMetadata{
		bitFields: make([]*bitFieldMetadata,
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

func (m *wordMetadata) marshal(reflection reflect.Value) (bytes []byte) {
	var (
		bitField       *bitFieldMetadata
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

func (m *wordMetadata) unmarshal(bytes []byte, reflection reflect.Value) {
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

type bitFieldMetadata struct {
	length uint
	offset uint64
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

func (m *bitFieldMetadata) marshal(reflection reflect.Value) (value uint64) {
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

	value = value & (1<<m.length - 1) << m.offset // XXX: mask if overflowing

	return
}

func (m *bitFieldMetadata) unmarshal(bytes []byte, reflection reflect.Value) {
	var (
		value uint64
	)

	value = binary.BigEndian.Uint64(bytes) >> m.offset & (1<<m.length - 1)

	switch m.kind {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough

	case reflect.Uint:
		reflection.SetUint(value)

	case reflect.Bool:
		switch value {
		case 1:
			reflection.SetBool(true)
		default:
			reflection.SetBool(false)
		}
	}

	return
}
