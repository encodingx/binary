package words

import (
	"fmt"
	"reflect"

	"github.com/joel-ling/go-bitfields/pkg/structs/fields"
	"github.com/joel-ling/go-bitfields/pkg/structs/errors"
)

const (
	bitsPerByte = 8
)

type Word interface {
	// A Word represents a sequence of bytes made up of bit Fields.
	// See package `fields` for the definition of a Field.

	LengthInBits() uint
	// Return the length of a Word in number of bits.

	LengthInBytes() uint
	// Return the length of a Word in number of bytes.

	NFields() uint
	// Return the number of Fields in a Word.

	Field(uint) fields.Field
	// Return a Field given its index in a Word.
}

type word struct {
	sliceOfFields []fields.Field
	lengthInBits  uint
}

func NewWordFromStructField(structField reflect.StructField, formatName string,
) (
	w *word, e error,
) {
	// Return a default implementation of the interface Word,
	// given a reflection of a struct field representing a Word in a Field,
	// and the name of the Format containing that Word.

	const (
		structFieldTagKey         = "word"
		structFieldTagValueFormat = "%d"

		wordLengthLimit = 64
	)

	var (
		field               fields.Field
		i                   int
		structFieldTagKeyOK bool
		structFieldTagValue string
		sumFieldLengths     uint
	)

	if structField.Type.Kind() != reflect.Struct {
		e = errors.NewWordIsNotStructError(
			formatName,
			structField.Name,
		)

		return
	}

	if len(structField.Tag) == 0 {
		e = errors.NewWordMissingStructTagError(
			formatName,
			structField.Name,
		)

		return
	}

	w = new(word)

	structFieldTagValue, structFieldTagKeyOK = structField.Tag.Lookup(
		structFieldTagKey,
	)

	_, e = fmt.Sscanf(
		structFieldTagValue,
		structFieldTagValueFormat,
		&w.lengthInBits,
	)
	if e != nil || !structFieldTagKeyOK {
		e = errors.NewWordStructTagMalformedError(
			formatName,
			structField.Name,
		)

		return
	}

	if w.lengthInBits > wordLengthLimit {
		e = errors.NewWordLengthExceedsLimitError(
			formatName,
			structField.Name,
			w.lengthInBits,
			wordLengthLimit,
		)

		return
	}

	if w.lengthInBits < 1 || w.lengthInBits%bitsPerByte != 0 {
		e = errors.NewWordLengthNotMultipleOfFactorError(
			formatName,
			structField.Name,
			w.lengthInBits,
			bitsPerByte,
		)

		return
	}

	if structField.Type.NumField() == 0 {
		e = errors.NewWordHasNoFieldsError(
			formatName,
			structField.Name,
		)

		return
	}

	w.sliceOfFields = make([]fields.Field,
		structField.Type.NumField(),
	)

	for i = structField.Type.NumField() - 1; i >= 0; i-- {
		field, sumFieldLengths, e = fields.NewFieldFromStructField(
			structField.Type.Field(i),
			structField.Name,
			formatName,
			sumFieldLengths,
		)
		if e != nil {
			return
		}

		w.sliceOfFields[i] = field
	}

	if sumFieldLengths != w.lengthInBits {
		e = errors.NewWordLengthNotSumOfFieldLengthsError(
			formatName,
			structField.Name,
			w.lengthInBits,
			sumFieldLengths,
		)

		return
	}

	return
}

func (w *word) LengthInBits() uint {
	// Implement interface Word.

	return w.lengthInBits
}

func (w *word) LengthInBytes() uint {
	// Implement interface Word.

	return w.lengthInBits / bitsPerByte
}

func (w *word) NFields() (nFields uint) {
	// Implement interface Word.

	nFields = uint(
		len(w.sliceOfFields),
	)

	return
}

func (w *word) Field(index uint) fields.Field {
	// Implement interface Word.

	return w.sliceOfFields[index]
}
