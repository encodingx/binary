package formats

import (
	"reflect"
	"testing"

	"github.com/joel-ling/go-bitfields/internal/structs/test"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	const (
		nCases = 13
	)

	var (
		testFormat14Bad = test.TestFormat14Bad(0)

		formats = []interface{}{
			&test.TestFormat0{},
			&test.TestFormat1Bad{},
			&test.TestFormat2Bad{},
			&test.TestFormat5Bad{},
			&test.TestFormat6Bad{},
			&test.TestFormat7Bad{},
			&test.TestFormat8Bad{},
			&test.TestFormat9Bad{},
			&test.TestFormat10Bad{},
			&test.TestFormat11Bad{},
			&test.TestFormat12Bad{},
			&test.TestFormat13Bad{},
			&testFormat14Bad,
		}

		lengthsInBits = []uint{
			96,
			96,
			96,
			96,
			96,
			96,
			64,
			96,
			88,
			136,
			0,
			96,
			0,
		}

		lengthsInBytes = []uint{
			12,
			12,
			12,
			12,
			12,
			12,
			8,
			12,
			11,
			11,
			0,
			12,
			0,
		}

		nWords = []uint{
			3,
			3,
			3,
			3,
			3,
			3,
			3,
			3,
			3,
			3,
			0,
			3,
			0,
		}

		expectError = []bool{
			false,
			true,
			true,
			true,
			true,
			true,
			true,
			true,
			true,
			true,
			true,
			true,
			true,
		}

		e error
		f Format
		i int
	)

	for i = 0; i < nCases; i++ {
		f, e = NewFormatFromType(
			reflect.TypeOf(formats[i]),
		)

		if expectError[i] {
			assert.NotNil(t, e,
				"The constructor of a Format "+
					"should return an error "+
					"when it is passed a bad struct.",
			)

			t.Log(e)

			continue

		} else {
			if e != nil {
				t.Error(e)
			}
		}

		assert.Equal(t,
			lengthsInBits[i], f.LengthInBits(),
			"The value returned by Format.LengthInBits() "+
				"should be equal to the expected value given the struct.",
		)

		assert.Equal(t,
			lengthsInBytes[i], f.LengthInBytes(),
			"The value returned by Format.LengthInBytes() "+
				"should be equal to the expected value given the struct.",
		)

		assert.Equal(t,
			nWords[i], f.NWords(),
			"The value returned by Format.NWords() "+
				"should be equal to the expected value given the struct.",
		)
	}
}
