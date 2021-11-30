package formats

import (
	"reflect"

	"github.com/joel-ling/go-bitfields/pkg/structs/errors"
	"github.com/joel-ling/go-bitfields/pkg/structs/words"
)

type Format interface {
	// A Format represents a binary message or file format made up of Words.
	// See package `words` for the defintion of a Word.

	LengthInBits() uint
	// Return the sum of the lengths of Words in the Format, in number of bits.

	LengthInBytes() uint
	// Return the sum of the lengths of Words in the Format, in number of bytes.

	NWords() uint
	// Return the number of Words in a Format.

	Word(uint) words.Word
	// Return a Word given its index in a Format.
}

type format struct {
	sliceOfWords []words.Word
	lengthInBits uint
}

func NewFormatFromType(typeReflection reflect.Type) (f *format, e error) {
	// Return a default implementation of the interface Format,
	// given a type reflection of a format struct.

	var (
		i    int
		word words.Word
	)

	if typeReflection.Kind() != reflect.Ptr {
		e = errors.NewInterfaceIsNotPointerError(
			typeReflection.String(),
		)

		return
	}

	typeReflection = typeReflection.Elem()

	if typeReflection.Kind() != reflect.Struct {
		e = errors.NewFormatIsNotStructError(
			typeReflection.String(),
		)

		return
	}

	if typeReflection.NumField() == 0 {
		e = errors.NewFormatHasNoWordsError(
			typeReflection.String(),
		)

		return
	}

	f = &format{
		sliceOfWords: make([]words.Word,
			typeReflection.NumField(),
		),
	}

	for i = 0; i < typeReflection.NumField(); i++ {
		word, e = words.NewWordFromStructField(
			typeReflection.Field(i),
			typeReflection.String(),
		)
		if e != nil {
			return
		}

		f.sliceOfWords[i] = word

		f.lengthInBits += word.LengthInBits()
	}

	return
}

func (f *format) LengthInBits() uint {
	// Implement interface Format.

	return f.lengthInBits
}

func (f *format) LengthInBytes() uint {
	// Implement interface Format.

	const (
		bitsPerByte = 8
	)

	return f.lengthInBits / bitsPerByte
}

func (f *format) NWords() (nWords uint) {
	// Implement interface Format.

	nWords = uint(
		len(f.sliceOfWords),
	)

	return
}

func (f *format) Word(index uint) words.Word {
	// Implement interface Format.

	return f.sliceOfWords[index]
}
