package errors

import (
	"fmt"
)

type wordLengthNotMultipleOfFactorError struct {
	wordError
	wordLength uint
	factor     uint
}

func NewWordLengthNotMultipleOfFactorError(formatName, wordName string,
	wordLength, factor uint,
) (
	e error,
) {
	e = &wordLengthNotMultipleOfFactorError{
		wordError: wordError{
			formatError: formatError{
				formatName: formatName,
			},
			wordName: wordName,
		},
		wordLength: wordLength,
		factor:     factor,
	}

	return
}

func (e wordLengthNotMultipleOfFactorError) Error() (message string) {
	const (
		format = "Length %d bits " +
			"of word %s " +
			"in format %s " +
			"is not a multiple of %d."
	)

	message = fmt.Sprintf(format,
		e.wordLength,
		e.wordName,
		e.formatName,
		e.factor,
	)

	return
}
