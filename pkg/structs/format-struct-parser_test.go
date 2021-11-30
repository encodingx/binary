package structs

import (
	"testing"

	"github.com/joel-ling/go-bitfields/pkg/structs/formats"
	"github.com/joel-ling/go-bitfields/pkg/structs/formats/test"
	"github.com/stretchr/testify/assert"
)

func TestFormatStructParser(t *testing.T) {
	const (
		nWords uint = 3
	)

	var (
		wordLengthsInBits = []uint{
			24,
			32,
			40,
		}

		wordLengthsInBytes = []uint{
			3,
			4,
			5,
		}

		nFields = []uint{
			3,
			4,
			5,
		}

		fieldLengthsInBits = [][]uint{
			[]uint{
				8,
				8,
				8,
			},
			[]uint{
				5,
				7,
				9,
				11,
			},
			[]uint{
				3,
				10,
				6,
				13,
				8,
			},
		}

		fieldOffsetsInBits = [][]uint{
			[]uint{
				16,
				8,
				0,
			},
			[]uint{
				27,
				20,
				11,
				0,
			},
			[]uint{
				37,
				27,
				21,
				8,
				0,
			},
		}

		e      error
		i      uint
		j      uint
		parser FormatStructParser

		format  formats.Format
		struct0 test.TestFormat0
	)

	parser = NewFormatStructParser()

	format, e = parser.ParseFormatStruct(&struct0)
	if e != nil {
		t.Error(e)
	}

	assert.Equal(t,
		nWords, format.NWords(),
		"Method NWords() of a Format parsed should "+
			"return the expected number of Words in the Format.",
	)

	for i = 0; i < nWords; i++ {
		assert.Equal(t,
			wordLengthsInBits[i], format.Word(i).LengthInBits(),
			"Chained methods Word(i).LengthInBits() of a Format parsed should "+
				"return the expected Word length in bits.",
		)

		assert.Equal(t,
			wordLengthsInBytes[i], format.Word(i).LengthInBytes(),
			"Chained methods Word(i).LengthInBytes() of a Format parsed "+
				"should return the expected Word length in bytes.",
		)

		assert.Equal(t,
			nFields[i], format.Word(i).NFields(),
			"Chained methods Word(i).NFields() of a Format parsed should "+
				"return the expected number of fields in the Word.",
		)

		for j = 0; j < format.Word(i).NFields(); j++ {
			assert.Equal(t,
				fieldLengthsInBits[i][j],
				format.Word(i).Field(j).LengthInBits(),
				"Chained methods Word(i).Field(j).LengthInBits "+
					"of a Format parsed "+
					"should return the expected Field length in bits.",
			)

			assert.Equal(t,
				fieldOffsetsInBits[i][j],
				format.Word(i).Field(j).OffsetInBits(),
				"Chained methods Word(i).Field(j).OffsetInBits "+
					"of a Format parsed "+
					"should return the expected Field offset in bits.",
			)
		}
	}
}
