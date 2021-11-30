package errors

import (
	"fmt"
)

type wordLengthNotSumOfFieldLengthsError struct {
	wordError
	wordLength        uint
	sumOfFieldLengths uint
}

func NewWordLengthNotSumOfFieldLengthsError(formatName, wordName string,
	wordLength, sumOfFieldLengths uint,
) (
	e error,
) {
	e = &wordLengthNotSumOfFieldLengthsError{
		wordError: wordError{
			formatError: formatError{
				formatName: formatName,
			},
			wordName: wordName,
		},
		wordLength:        wordLength,
		sumOfFieldLengths: sumOfFieldLengths,
	}

	return
}

func (e wordLengthNotSumOfFieldLengthsError) Error() (message string) {
	const (
		format = "Length %d bits declared in struct tag " +
			"of word %s " +
			"in format %s " +
			"is not equal to the sum of the lengths of its fields, %d bits."
	)

	message = fmt.Sprintf(format,
		e.wordLength,
		e.wordName,
		e.formatName,
		e.sumOfFieldLengths,
	)

	return
}
