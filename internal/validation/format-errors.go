package validation

import (
	"fmt"
)

type FormatError interface {
	FunctionError
	SetFormatName(string)
}

type DefaultFormatError struct {
	DefaultFunctionError
	formatName string
}

func (e *DefaultFormatError) SetFormatName(formatName string) {
	e.formatName = formatName

	return
}

type formatWithNoWordsError struct {
	DefaultFormatError
}

func NewFormatWithNoWordsError() *formatWithNoWordsError {
	return new(formatWithNoWordsError)
}

func (e *formatWithNoWordsError) Error() string {
	const (
		format = "" +
			"A format-struct should nest exported word-structs. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"that has no words."
	)

	return fmt.Sprintf(format, e.functionName, e.formatName)
}

type lengthOfByteSliceNotEqualToFormatLengthError struct {
	DefaultFormatError
	formatLength    uint
	byteSliceLength uint
}

func NewLengthOfByteSliceNotEqualToFormatLengthError(
	formatLength, byteSliceLength uint,
) (
	e *lengthOfByteSliceNotEqualToFormatLengthError,
) {
	e = &lengthOfByteSliceNotEqualToFormatLengthError{
		formatLength:    formatLength,
		byteSliceLength: byteSliceLength,
	}

	return
}

func (e *lengthOfByteSliceNotEqualToFormatLengthError) Error() string {
	const (
		format = "" +
			"A byte slice into which a format-struct would be unmarshalled " +
			"should be of length equal to the sum of lengths of words " +
			"in the format represented by the struct. " +
			"Argument to Unmarshal points to a format-struct \"%s\" " +
			"of length %d bits " +
			"not equal to the length of the byte slice, %d bits."
	)

	return fmt.Sprintf(format, e.formatName, e.formatLength, e.byteSliceLength)
}
