package validation

import (
	"fmt"
)

type WordError interface {
	FormatError
	SetWordName(string)
}

type DefaultWordError struct {
	DefaultFormatError
	wordName string
}

func (e *DefaultWordError) SetWordName(wordName string) {
	e.wordName = wordName

	return
}

type wordNotStructError struct {
	DefaultWordError
}

func NewWordNotStructError() *wordNotStructError {
	return new(wordNotStructError)
}

func (e *wordNotStructError) Error() string {
	const (
		format = "" +
			"A format-struct should nest exported word-structs. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"that is not a struct."
	)

	return fmt.Sprintf(format, e.functionName, e.formatName, e.wordName)
}

type wordOfIncompatibleLengthError struct {
	DefaultWordError
	wordLength uint
}

func NewWordOfIncompatibleLengthError(wordLength uint) (
	e *wordOfIncompatibleLengthError,
) {
	e = &wordOfIncompatibleLengthError{
		wordLength: wordLength,
	}

	return
}

func (e *wordOfIncompatibleLengthError) Error() (s string) {
	const (
		format = "" +
			"The length of a word should be a multiple of eight " +
			"in the range [8, 64]. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"that has a word \"%s\" " +
			"of length %d not in {8, 16, 24, ... 64}."
	)

	s = fmt.Sprintf(format,
		e.functionName, e.formatName, e.wordName,
		e.wordLength,
	)

	return
}

type wordOfLengthNotEqualToSumOfLengthsOfBitFieldsError struct {
	DefaultWordError
	wordLength        uint
	bitFieldLengthSum uint
}

func NewWordOfLengthNotEqualToSumOfLengthsOfBitFieldsError(
	wordLength, bitFieldLengthSum uint,
) (
	e *wordOfLengthNotEqualToSumOfLengthsOfBitFieldsError,
) {
	e = &wordOfLengthNotEqualToSumOfLengthsOfBitFieldsError{
		wordLength:        wordLength,
		bitFieldLengthSum: bitFieldLengthSum,
	}

	return
}

func (e *wordOfLengthNotEqualToSumOfLengthsOfBitFieldsError) Error() (s string) {
	const (
		format = "" +
			"The length of a word " +
			"should be equal to the sum of lengths of its bit fields. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"that has a word \"%s\" " +
			"of length %d " +
			"not equal to the sum of the lengths of its bit fields, %d."
	)

	s = fmt.Sprintf(format,
		e.functionName, e.formatName, e.wordName,
		e.wordLength, e.bitFieldLengthSum,
	)

	return
}

type wordWithMalformedTagError struct {
	DefaultWordError
}

func NewWordWithMalformedTagError() *wordWithMalformedTagError {
	return new(wordWithMalformedTagError)
}

func (e *wordWithMalformedTagError) Error() string {
	const (
		format = "" +
			"A format-struct should nest exported word-structs " +
			"tagged with a key \"word\" and a value " +
			"indicating the length of a word in number of bits " +
			"(e.g. `word:\"32\"`). " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"with a malformed struct tag."
	)

	return fmt.Sprintf(format, e.functionName, e.formatName, e.wordName)
}

type wordWithNoBitFieldsError struct {
	DefaultWordError
}

func NewWordWithNoBitFieldsError() *wordWithNoBitFieldsError {
	return new(wordWithNoBitFieldsError)
}

func (e *wordWithNoBitFieldsError) Error() string {
	const (
		format = "" +
			"A word-struct should have exported fields " +
			"representing bit fields. " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"that has no bit fields."
	)

	return fmt.Sprintf(format, e.functionName, e.formatName, e.wordName)
}

type wordWithNoStructTagError struct {
	DefaultWordError
}

func NewWordWithNoStructTagError() *wordWithNoStructTagError {
	return new(wordWithNoStructTagError)
}

func (e *wordWithNoStructTagError) Error() string {
	const (
		format = "" +
			"A format-struct should nest exported word-structs " +
			"tagged with a key \"word\" and a value " +
			"indicating the length of a word in number of bits " +
			"(e.g. `word:\"32\"`). " +
			"Argument to %s points to a format-struct \"%s\" " +
			"nesting a word-struct \"%s\" " +
			"with no struct tag."
	)

	return fmt.Sprintf(format, e.functionName, e.formatName, e.wordName)
}
