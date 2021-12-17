package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	wordName = "Word"
)

func TestWordNotStructError(t *testing.T) {
	const (
		errorMessage = "" +
			"A format-struct should nest exported word-structs. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"that is not a struct."
	)

	var (
		e WordError
	)

	e = NewWordNotStructError()

	e.SetFunctionName(functionName)

	e.SetFormatName(formatName)

	e.SetWordName(wordName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}

func TestWordOfIncompatibleLengthError(t *testing.T) {
	const (
		wordLength = 36

		errorMessage = "" +
			"The length of a word should be a multiple of eight " +
			"in the range [8, 64]. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"that has a word \"Word\" " +
			"of length 36 not in {8, 16, 24, ... 64}."
	)

	var (
		e WordError
	)

	e = NewWordOfIncompatibleLengthError(wordLength)

	e.SetFunctionName(functionName)

	e.SetFormatName(formatName)

	e.SetWordName(wordName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}

func TestWordOfLengthNotEqualToSumOfLengthsOfBitFieldsError(t *testing.T) {
	const (
		bitFieldLengthSum = 31
		wordLength        = 32

		errorMessage = "" +
			"The length of a word " +
			"should be equal to the sum of lengths of its bit fields. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"that has a word \"Word\" " +
			"of length 32 " +
			"not equal to the sum of the lengths of its bit fields, 31."
	)

	var (
		e WordError
	)

	e = NewWordOfLengthNotEqualToSumOfLengthsOfBitFieldsError(
		wordLength, bitFieldLengthSum,
	)

	e.SetFunctionName(functionName)

	e.SetFormatName(formatName)

	e.SetWordName(wordName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}

func TestWordWithMalformedTagError(t *testing.T) {
	const (
		errorMessage = "" +
			"A format-struct should nest exported word-structs " +
			"tagged with a key \"word\" and a value " +
			"indicating the length of a word in number of bits " +
			"(e.g. `word:\"32\"`). " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"with a malformed struct tag."
	)

	var (
		e WordError
	)

	e = NewWordWithMalformedTagError()

	e.SetFunctionName(functionName)

	e.SetFormatName(formatName)

	e.SetWordName(wordName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}

func TestWordWithNoBitFieldsError(t *testing.T) {
	const (
		errorMessage = "" +
			"A word-struct should have exported fields " +
			"representing bit fields. " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has no bit fields."
	)

	var (
		e WordError
	)

	e = NewWordWithNoBitFieldsError()

	e.SetFunctionName(functionName)

	e.SetFormatName(formatName)

	e.SetWordName(wordName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}

func TestWordWithNoStructTagError(t *testing.T) {
	const (
		errorMessage = "" +
			"A format-struct should nest exported word-structs " +
			"tagged with a key \"word\" and a value " +
			"indicating the length of a word in number of bits " +
			"(e.g. `word:\"32\"`). " +
			"Argument to Marshal points to a format-struct \"Format\" " +
			"nesting a word-struct \"Word\" " +
			"with no struct tag."
	)

	var (
		e WordError
	)

	e = NewWordWithNoStructTagError()

	e.SetFunctionName(functionName)

	e.SetFormatName(formatName)

	e.SetWordName(wordName)

	assert.Equal(t,
		errorMessage, e.Error(),
	)
}
