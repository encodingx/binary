package errors

import (
	"fmt"
)

type wordLengthExceedsLimitError struct {
	wordError
	wordLength      uint
	wordLengthLimit uint
}

func NewWordLengthExceedsLimitError(formatName, wordName string,
	wordLength, wordLengthLimit uint,
) (
	e error,
) {
	e = &wordLengthExceedsLimitError{
		wordError: wordError{
			formatName: formatName,
			wordName:   wordName,
		},
		wordLength:      wordLength,
		wordLengthLimit: wordLengthLimit,
	}

	return
}

func (e wordLengthExceedsLimitError) Error() (message string) {
	const (
		format = "Length %d bits " +
			"of word %s " +
			"in format %s " +
			"exceeds the word length limit of %d bits."
	)

	message = fmt.Sprintf(format,
		e.wordLength,
		e.wordName,
		e.formatName,
		e.wordLengthLimit,
	)

	return
}
