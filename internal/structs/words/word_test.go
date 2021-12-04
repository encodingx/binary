package words

import (
	"reflect"
	"testing"

	"github.com/encodingx/binary/internal/structs/test"
	"github.com/stretchr/testify/assert"
)

func TestWord(t *testing.T) {
	const (
		nCases = 10
	)

	var (
		formats = []interface{}{
			test.TestFormat0{},
			test.TestFormat1Bad{},
			test.TestFormat2Bad{},
			test.TestFormat5Bad{},
			test.TestFormat6Bad{},
			test.TestFormat7Bad{},
			test.TestFormat9Bad{},
			test.TestFormat10Bad{},
			test.TestFormat11Bad{},
			test.TestFormat13Bad{},
		}

		structFieldIndices = []int{
			0,
			0,
			1,
			1,
			2,
			0,
			2,
			0,
			1,
			2,
		}

		lengthsInBits = []uint{
			24,
			24,
			32,
			32,
			40,
			20,
			40,
			16,
			72,
			40,
		}

		lengthsInBytes = []uint{
			3,
			3,
			4,
			4,
			5,
			2,
			5,
			2,
			9,
			5,
		}

		nFields = []uint{
			3,
			3,
			4,
			4,
			6,
			3,
			5,
			3,
			4,
			5,
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
		}

		e              error
		i              int
		typeReflection reflect.Type
		w              Word
	)

	for i = 0; i < nCases; i++ {
		typeReflection = reflect.TypeOf(formats[i])

		w, e = NewWordFromStructField(
			typeReflection.Field(structFieldIndices[i]),
			typeReflection.Name(),
		)

		if expectError[i] {
			assert.NotNil(t, e,
				"The constructor of a Word "+
					"should return an error "+
					"when it is passed a bad struct field.",
			)

			t.Log(e)

			continue

		} else {
			if e != nil {
				t.Error(e)
			}
		}

		assert.Equal(t,
			lengthsInBits[i], w.LengthInBits(),
			"The value returned by Field.LengthInBits() "+
				"should be equal to the relevant value in the struct tag.",
		)

		assert.Equal(t,
			lengthsInBytes[i], w.LengthInBytes(),
			"The value returned by Field.LengthInBytes() "+
				"should be equal to the expected value given the struct tag.",
		)

		assert.Equal(t,
			nFields[i], w.NFields(),
			"The value returned by Field.NFields() "+
				"should be equal to the expected value given the struct.",
		)
	}
}
