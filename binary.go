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
			e = fmt.Errorf(marshalError, e)
		}

		return
	}()

	reflection, e = structReflectionFromInterface(iface, functionName)
	if e != nil {
		return
	}

	_, e = validateFormatReflection(reflection, functionName)
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
			e = fmt.Errorf(unmarshalError, e)
		}

		return
	}()

	reflection, e = structReflectionFromInterface(iface, functionName)
	if e != nil {
		return
	}

	formatLength, e = validateFormatReflection(reflection, functionName)
	if e != nil {
		return
	}

	byteSliceLength = uint(
		len(bytes) * bitsPerByte,
	)

	if byteSliceLength != formatLength {
		e = validation.NewLengthOfByteSliceNotEqualToFormatLengthError(
			reflection.String(),
			formatLength,
			byteSliceLength,
		)

		return
	}

	return
}

func structReflectionFromInterface(
	iface interface{}, functionName string,
) (
	reflection reflect.Type, e error,
) {
	reflection = reflect.TypeOf(iface)

	if reflection.Kind() != reflect.Ptr {
		e = validation.NewNonPointerError(functionName)

		return
	}

	reflection = reflection.Elem()

	if reflection.Kind() != reflect.Struct {
		e = validation.NewPointerToNonStructVariableError(functionName)

		return
	}

	return
}

func validateFormatReflection(reflection reflect.Type, functionName string) (
	formatLength uint, e error,
) {
	var (
		i          int
		wordLength uint
	)

	if reflection.NumField() == 0 {
		e = validation.NewFormatWithNoWordsError(
			functionName,
			reflection.String(),
		)

		return
	}

	for i = 0; i < reflection.NumField(); i++ {
		wordLength, e = validateWordReflection(
			reflection.Field(i),
			functionName,
			reflection.String(),
		)
		if e != nil {
			return
		}

		formatLength += wordLength
	}

	return
}

func validateWordReflection(
	reflection reflect.StructField, functionName, formatName string,
) (
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

	if reflection.Type.Kind() != reflect.Struct {
		e = validation.NewWordNotStructError(
			functionName,
			formatName,
			reflection.Name,
		)

		return
	}

	if len(reflection.Tag) == 0 {
		e = validation.NewWordWithNoStructTagError(
			functionName,
			formatName,
			reflection.Name,
		)

		return
	}

	_, e = fmt.Sscanf(
		reflection.Tag.Get(tagKey),
		tagValueFormat,
		&wordLength,
	)
	if e != nil {
		e = validation.NewWordWithMalformedTagError(
			functionName,
			formatName,
			reflection.Name,
		)

		return
	}

	wordLengthOK = wordLength%wordLengthFactor == 0
	wordLengthOK = wordLengthOK && wordLength >= wordLengthLowerLimit
	wordLengthOK = wordLengthOK && wordLength <= wordLengthUpperLimit

	if !wordLengthOK {
		e = validation.NewWordOfIncompatibleLengthError(
			functionName,
			formatName,
			reflection.Name,
			wordLength,
		)

		return
	}

	if reflection.Type.NumField() == 0 {
		e = validation.NewWordWithNoBitFieldsError(
			functionName,
			formatName,
			reflection.Name,
		)

		return
	}

	for i = 0; i < reflection.Type.NumField(); i++ {
		bitFieldLength, e = validateBitFieldReflection(
			reflection.Type.Field(i),
			functionName,
			formatName,
			reflection.Name,
		)
		if e != nil {
			return
		}

		bitFieldLengthSum += bitFieldLength
	}

	if bitFieldLengthSum != wordLength {
		e = validation.NewWordOfLengthNotEqualToSumOfLengthsOfBitFieldsError(
			functionName,
			formatName,
			reflection.Name,
			wordLength,
			bitFieldLengthSum,
		)

		return
	}

	return
}

func validateBitFieldReflection(
	reflection reflect.StructField, functionName, formatName, wordName string,
) (
	bitFieldLength uint, e error,
) {
	const (
		tagKey         = "bitfield"
		tagValueFormat = "%d"
	)

	var (
		bitFieldLengthCap uint
	)

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
			functionName,
			formatName,
			wordName,
			reflection.Name,
			reflection.Type.String(),
		)

		return
	}

	if len(reflection.Tag) == 0 {
		e = validation.NewBitFieldWithNoStructTagError(
			functionName,
			formatName,
			wordName,
			reflection.Name,
		)

		return
	}

	_, e = fmt.Sscanf(
		reflection.Tag.Get(tagKey),
		tagValueFormat,
		&bitFieldLength,
	)
	if e != nil {
		e = validation.NewBitFieldWithMalformedTagError(
			functionName,
			formatName,
			wordName,
			reflection.Name,
		)

		return
	}

	if bitFieldLength > bitFieldLengthCap {
		e = validation.NewBitFieldOfLengthOverflowingTypeError(
			functionName,
			formatName,
			wordName,
			reflection.Name,
			bitFieldLength,
			reflection.Type.String(),
		)
	}

	return
}
