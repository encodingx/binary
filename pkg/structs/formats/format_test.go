package formats

import (
	"reflect"
	"testing"

	"github.com/joel-ling/go-bitfields/pkg/structs/formats/test"
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

		nWords = []uint{
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
			nWords[i], f.NWords(),
			"The value returned by Format.NWords() "+
				"should be equal to the expected value given the struct.",
		)
	}
}
