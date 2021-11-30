package errors

import (
	"fmt"
)

type wordError struct {
	formatName string
	wordName   string
}

type wordLengthNotMultipleOf8Error struct {
	wordError
	wordLength uint
}

func NewWordLengthNotMultipleOf8Error(formatName, wordName string,
	wordLength uint,
) (
	e error,
) {
	e = &wordLengthNotMultipleOf8Error{
		wordError: wordError{
			formatName: formatName,
			wordName:   wordName,
		},
		wordLength: wordLength,
	}

	return
}

func (e wordLengthNotMultipleOf8Error) Error() (message string) {
	const (
		format = "Length %d bits " +
			"of word %s " +
			"in format %s " +
			"is not a multiple of eight."
	)

	message = fmt.Sprintf(format,
		e.wordLength,
		e.wordName,
		e.formatName,
	)

	return
}
