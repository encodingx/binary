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
	formatLengthInBytes uint
	byteSliceLength     uint
}

func NewLengthOfByteSliceNotEqualToFormatLengthError(
	formatLengthInBytes, byteSliceLength uint,
) (
	e *lengthOfByteSliceNotEqualToFormatLengthError,
) {
	e = &lengthOfByteSliceNotEqualToFormatLengthError{
		formatLengthInBytes: formatLengthInBytes,
		byteSliceLength:     byteSliceLength,
	}

	return
}

func (e *lengthOfByteSliceNotEqualToFormatLengthError) Error() (s string) {
	const (
		format = "" +
			"A byte slice into which a format-struct would be unmarshalled " +
			"should be of length equal to the sum of lengths of words " +
			"in the format represented by the struct. " +
			"Argument to Unmarshal points to a format-struct \"%s\" " +
			"of length %d byte(s) " +
			"not equal to the length of the byte slice, %d byte(s)."
	)

	s = fmt.Sprintf(format,
		e.formatName,
		e.formatLengthInBytes,
		e.byteSliceLength,
	)

	return
}
