package validation

import (
	"fmt"
)

type BitFieldError interface {
	WordError
	SetBitFieldName(string)
}

type DefaultBitFieldError struct {
	DefaultWordError
	bitFieldName string
}

func (e *DefaultBitFieldError) SetBitFieldName(bitFieldName string) {
	e.bitFieldName = bitFieldName

	return
}

type bitFieldOfUnsupportedTypeError struct {
	DefaultBitFieldError
	bitFieldType string
}

func NewBitFieldOfUnsupportedTypeError(bitFieldType string) (
	e *bitFieldOfUnsupportedTypeError,
) {
	e = &bitFieldOfUnsupportedTypeError{
		bitFieldType: bitFieldType,
	}

	return
}

func (e *bitFieldOfUnsupportedTypeError) Error() (s string) {
	const (
		format = "" +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"of type uintN or bool. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"that has a bit field \"%s\" " +
			"of unsupported type \"%s\"."
	)

	s = fmt.Sprintf(format,
		e.functionName,
		e.formatName,
		e.wordName,
		e.bitFieldName,
		e.bitFieldType,
	)

	return
}

type bitFieldOfLengthOverflowingTypeError struct {
	DefaultBitFieldError
	bitFieldLength uint
	bitFieldType   string
}

func NewBitFieldOfLengthOverflowingTypeError(
	bitFieldLength uint, bitFieldType string,
) (
	e *bitFieldOfLengthOverflowingTypeError,
) {
	e = &bitFieldOfLengthOverflowingTypeError{
		bitFieldLength: bitFieldLength,
		bitFieldType:   bitFieldType,
	}

	return
}

func (e *bitFieldOfLengthOverflowingTypeError) Error() (s string) {
	const (
		format = "" +
			"The number of unique values a bit field can contain " +
			"must not exceed the size of its type. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"that has a bit field \"%s\" " +
			"of length %d exceeding the size of type \"%s\"."
	)

	s = fmt.Sprintf(format,
		e.functionName, e.formatName, e.wordName, e.bitFieldName,
		e.bitFieldLength, e.bitFieldType,
	)

	return
}

type bitFieldWithMalformedTagError struct {
	DefaultBitFieldError
}

func NewBitFieldWithMalformedTagError() *bitFieldWithMalformedTagError {
	return new(bitFieldWithMalformedTagError)
}

func (e *bitFieldWithMalformedTagError) Error() (s string) {
	const (
		format = "" +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"tagged with a key \"bitfield\" and a value " +
			"indicating the length of the bit field in number of bits " +
			"(e.g. `bitfield:\"1\"`). " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"that has a bit field \"%s\" " +
			"with a malformed struct tag."
	)

	s = fmt.Sprintf(format,
		e.functionName, e.formatName, e.wordName, e.bitFieldName,
	)

	return
}

type bitFieldWithNoStructTagError struct {
	DefaultBitFieldError
}

func NewBitFieldWithNoStructTagError() *bitFieldWithNoStructTagError {
	return new(bitFieldWithNoStructTagError)
}

func (e *bitFieldWithNoStructTagError) Error() (s string) {
	const (
		format = "" +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"tagged with a key \"bitfield\" and a value " +
			"indicating the length of the bit field in number of bits " +
			"(e.g. `bitfield:\"1\"`). " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"that has a bit field \"%s\" " +
			"with no struct tag."
	)

	s = fmt.Sprintf(format,
		e.functionName, e.formatName, e.wordName, e.bitFieldName,
	)

	return
}
