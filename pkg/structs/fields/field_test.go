package fields

import (
	"reflect"
	"testing"

	"github.com/joel-ling/go-bitfields/pkg/structs/formats/test"
	"github.com/stretchr/testify/assert"
)

func TestField(t *testing.T) {
	const (
		formatName = "formatName"
		nCases     = 8
	)

	var (
		words = []interface{}{
			test.TestWord11Bad{},
			test.TestWord0{},
			test.TestWord7Bad{},
			test.TestWord13Bad{},
			test.TestWord12Bad{},
			test.TestWord3Bad{},
			test.TestWord4Bad{},
			test.TestWord8Bad{},
		}

		structFieldIndices = []int{
			0,
			0,
			1,
			4,
			2,
			0,
			1,
			2,
		}

		sumsFieldLengthsBefore = []uint{
			16,
			16,
			20,
			0,
			11,
			16,
			20,
			27,
		}

		sumsFieldLengthsAfter = []uint{
			24,
			24,
			27,
			8,
			20,
			24,
			27,
			27,
		}

		lengthsInBits = []uint{
			8,
			8,
			7,
			8,
			9,
			6,
			9,
			6,
		}

		offsetsInBits = []uint{
			16,
			16,
			20,
			0,
			11,
			18,
			18,
			21,
		}

		expectError = []bool{
			true,
			false,
			true,
			true,
			true,
			true,
			true,
			true,
		}

		e               error
		i               int
		f               Field
		typeReflection  reflect.Type
		sumFieldLengths uint
	)

	for i = 0; i < nCases; i++ {
		typeReflection = reflect.TypeOf(words[i])

		f, sumFieldLengths, e = NewFieldFromStructField(
			typeReflection.Field(structFieldIndices[i]),
			typeReflection.Name(),
			formatName,
			sumsFieldLengthsBefore[i],
		)

		if expectError[i] {
			assert.NotNil(t, e,
				"The constructor of a Field "+
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
			lengthsInBits[i], f.LengthInBits(),
			"The value returned by Field.LengthInBits() "+
				"should be equal to the relevant value in the struct tag.",
		)

		assert.Equal(t,
			offsetsInBits[i], f.OffsetInBits(),
			"The value returned by Field.OffsetInBits() "+
				"should be equal to the relevant value in the struct tag.",
		)

		assert.Equal(t,
			sumsFieldLengthsAfter[i], sumFieldLengths,
			"The cumulative sum of Field lengths "+
				"should be equal to the sum of the previous cumulative sum "+
				"and the length of the current Field.",
		)
	}
}
