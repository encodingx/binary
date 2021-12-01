package codecs

import (
	"testing"

	"github.com/joel-ling/go-bitfields/pkg/structs/test"
	"github.com/stretchr/testify/assert"
)

func TestV1Codec(t *testing.T) {
	const (
		testWord0TestField0 = 0b11111111
		testWord0TestField1 = 0
		testWord0TestField2 = 0b11111111

		testWord1TestField0 = 0
		testWord1TestField1 = 0b1111111
		testWord1TestField2 = 0
		testWord1TestField3 = 0b11111111111

		testWord2TestField0 = 0
		testWord2TestField1 = 0b1111111111
		testWord2TestField2 = 0
		testWord2TestField3 = 0b1111111111111
		testWord2TestField4 = 0
	)

	var (
		binary  []byte
		bytes   []byte
		codec   Codec
		e       error
		struct0 test.TestFormat0
		struct1 test.TestFormat0
		word0   test.TestWord0
		word1   test.TestWord1
		word2   test.TestWord2
	)

	word0 = test.TestWord0{
		TestField0: testWord0TestField0,
		TestField1: testWord0TestField1,
		TestField2: testWord0TestField2,
	}

	word1 = test.TestWord1{
		TestField0: testWord1TestField0,
		TestField1: testWord1TestField1,
		TestField2: testWord1TestField2,
		TestField3: testWord1TestField3,
	}

	word2 = test.TestWord2{
		TestField0: testWord2TestField0,
		TestField1: testWord2TestField1,
		TestField2: testWord2TestField2,
		TestField3: testWord2TestField3,
		TestField4: testWord2TestField4,
	}

	struct0 = test.TestFormat0{
		word0,
		word1,
		word2,
	}

	binary = []byte{
		0b11111111, 0b00000000, 0b11111111, // Word 0
		0b00000111, 0b11110000, 0b00000111, 0b11111111, // Word 1
		0b00011111, 0b11111000, 0b00011111, 0b11111111, 0b00000000, // Word 2
	}

	codec = NewV1Codec()

	bytes, e = codec.Marshal(&struct0)
	if e != nil {
		t.Error(e)
	}

	assert.Equal(t,
		binary, bytes,
		"The slice of bytes returned by method Codec.Marshal() "+
			"should match the expected given the struct field values.",
	)

	e = codec.Unmarshal(binary, &struct1)
	if e != nil {
		t.Error(e)
	}

	assert.Equal(t,
		struct0, struct1,
		"A struct marshalled by method Codec.Marshal() and "+
			"then unmarshalled by method Codec.Unmarshal() "+
			"should by the same as the original struct.",
	)
}
