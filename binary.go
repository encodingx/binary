package binary

import (
	"fmt"
	"reflect"

	"github.com/encodingx/binary/internal/validation"
)

func Marshal(iface interface{}) (bytes []byte, e error) {
	const (
		functionName = "Marshal"
	)

	var (
		reflection reflect.Type
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

	reflection, e = structReflectionFromInterface(iface)
	if e != nil {
		return
	}

	_, e = validateFormatReflection(reflection)
	if e != nil {
		return
	}

	return
}

func Unmarshal(bytes []byte, iface interface{}) (e error) {
	const (
		bitsPerByte  = 8
		functionName = "Unmarshal"
	)

	var (
		byteSliceLength uint
		formatLength    uint
		reflection      reflect.Type
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

	reflection, e = structReflectionFromInterface(iface)
	if e != nil {
		return
	}

	formatLength, e = validateFormatReflection(reflection)
	if e != nil {
		return
	}

	byteSliceLength = uint(
		len(bytes) * bitsPerByte,
	)

	if byteSliceLength != formatLength {
		e = validation.NewLengthOfByteSliceNotEqualToFormatLengthError(
			formatLength,
			byteSliceLength,
		)

        e.(validation.FormatError).SetFormatName(
            reflection.String(),
        )

		return
	}

	return
}

func structReflectionFromInterface(iface interface{}) (
	reflection reflect.Type, e error,
) {
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

	return
}

func validateFormatReflection(reflection reflect.Type) (
	formatLength uint, e error,
) {
	var (
		i          int
		wordLength uint
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

	for i = 0; i < reflection.NumField(); i++ {
		wordLength, e = validateWordReflection(
			reflection.Field(i),
		)
		if e != nil {
			return
		}

		formatLength += wordLength
	}

	return
}

func validateWordReflection(reflection reflect.StructField) (
	wordLength uint, e error,
) {
	const (
		tagKey         = "word"
		tagValueFormat = "%d"

		wordLengthFactor     = 8
		wordLengthLowerLimit = 8
		wordLengthUpperLimit = 64
	)

	var (
		bitFieldLength    uint
		bitFieldLengthSum uint
		wordLengthOK      bool

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

	for i = 0; i < reflection.Type.NumField(); i++ {
		bitFieldLength, e = validateBitFieldReflection(
			reflection.Type.Field(i),
		)
		if e != nil {
			return
		}

		bitFieldLengthSum += bitFieldLength
	}

	if bitFieldLengthSum != wordLength {
		e = validation.NewWordOfLengthNotEqualToSumOfLengthsOfBitFieldsError(
			wordLength,
			bitFieldLengthSum,
		)

		return
	}

	return
}

func validateBitFieldReflection(reflection reflect.StructField) (
	bitFieldLength uint, e error,
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

	if len(reflection.Tag) == 0 {
		e = validation.NewBitFieldWithNoStructTagError()

		return
	}

	_, e = fmt.Sscanf(
		reflection.Tag.Get(tagKey),
		tagValueFormat,
		&bitFieldLength,
	)
	if e != nil {
		e = validation.NewBitFieldWithMalformedTagError()

		return
	}

	if bitFieldLength > bitFieldLengthCap {
		e = validation.NewBitFieldOfLengthOverflowingTypeError(
			bitFieldLength,
			reflection.Type.String(),
		)
	}

	return
}
