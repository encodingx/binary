package binary

import (
	"fmt"
	"testing"

	"github.com/encodingx/binary/pkg/rfc791"
	"github.com/encodingx/binary/pkg/rfc791/v1p1"
	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	var (
		bytes []byte
		e     error
	)

	bytes, e = Marshal(&internetHeaderStruct)

	assert.Nil(t, e)

	assert.Equal(t,
		internetHeaderBytes, bytes,
	)
}

func TestMarshalV1p1(t *testing.T) {
	// Test backwards-compatibility with v1.1

	var (
		bytes []byte
		e     error
	)

	bytes, e = Marshal(&internetHeaderStructV1p1)

	assert.Nil(t, e)

	assert.Equal(t,
		internetHeaderBytes, bytes,
	)
}

func BenchmarkMarshal(b *testing.B) {
	var (
		e error
		i int
	)

	for i = 0; i < b.N; i++ {
		_, e = Marshal(&internetHeaderStruct)
		if e != nil {
			b.Error(e)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	var (
		e error
	)

	e = Unmarshal(internetHeaderBytes, &internetHeaderStruct1)

	assert.Nil(t, e)

	assert.Equal(t,
		internetHeaderStruct, internetHeaderStruct1,
	)
}

func TestUnmarshalV1p1(t *testing.T) {
	// Test backwards-compatibility with v1.1

	var (
		e error
	)

	e = Unmarshal(internetHeaderBytes, &internetHeaderStruct1V1p1)

	assert.Nil(t, e)

	assert.Equal(t,
		internetHeaderStructV1p1, internetHeaderStruct1V1p1,
	)
}

func BenchmarkUnmarshal(b *testing.B) {
	var (
		e error
		i int
	)

	for i = 0; i < b.N; i++ {
		e = Unmarshal(internetHeaderBytes, &internetHeaderStruct1)
		if e != nil {
			b.Error(e)
		}
	}
}

func TestShouldReturnErrorGivenNonPointer(t *testing.T) {
	const (
		errorMessage = "%[1]s error: " +
			"Argument to %[1]s should be a pointer to a format-struct. " +
			"Argument to %[1]s is not a pointer."
	)

	testShouldReturnErrorGiven(t,
		internetHeaderStruct,
		errorMessage,
	)
}

func TestShouldReturnErrorGivenPointerToNonStructVariable(t *testing.T) {
	const (
		errorMessage = "%[1]s error: " +
			"Argument to %[1]s should be a pointer to a format-struct. " +
			"Argument to %[1]s does not point to a struct variable."
	)

	testShouldReturnErrorGiven(t,
		&map[string]int{},
		errorMessage,
	)
}

func TestShouldReturnErrorGivenFormatWithNoWords(t *testing.T) {
	const (
		errorMessage = "%[1]s error: " +
			"A format-struct should nest exported word-structs. " +
			"Argument to %[1]s points to a format-struct \"binary.Format\" " +
			"that has no words."
	)

	type (
		Format struct{}
	)

	testShouldReturnErrorGiven(t,
		&Format{},
		errorMessage,
	)
}

func TestShouldReturnErrorGivenWordNotStruct(
	t *testing.T,
) {
	const (
		errorMessage = "%[1]s error: " +
			"A format-struct should nest exported word-structs. " +
			"Argument to %[1]s points to a format-struct \"binary.Format\" " +
			"nesting a word-struct \"Word\" " +
			"that is not a struct."
	)

	type (
		Format struct {
			Word map[string]int `word:"32"`
		}
	)

	testShouldReturnErrorGiven(t,
		&Format{},
		errorMessage,
	)
}

func TestShouldReturnErrorGivenWordWithNoStructTag(
	t *testing.T,
) {
	const (
		errorMessage = "%[1]s error: " +
			"A format-struct should nest exported word-structs " +
			"tagged with a key \"word\" and a value " +
			"indicating the length of a word in number of bits " +
			"(e.g. `word:\"32\"`). " +
			"Argument to %[1]s points to a format-struct \"binary.Format\" " +
			"nesting a word-struct \"Word\" " +
			"with no struct tag."
	)

	type (
		Word struct {
			BitField uint `bitfield:"32"`
		}

		Format struct {
			Word
		}
	)

	testShouldReturnErrorGiven(t,
		&Format{},
		errorMessage,
	)
}

func TestShouldReturnErrorGivenWordWithMalformedTag(
	t *testing.T,
) {
	const (
		errorMessage = "%[1]s error: " +
			"A format-struct should nest exported word-structs " +
			"tagged with a key \"word\" and a value " +
			"indicating the length of a word in number of bits " +
			"(e.g. `word:\"32\"`). " +
			"Argument to %[1]s points to a format-struct \"binary.Format\" " +
			"nesting a word-struct \"Word\" " +
			"with a malformed struct tag."
	)

	type (
		Word struct {
			BitField uint `bitfield:"32"`
		}

		Format struct {
			Word `worm:"32"`
		}
	)

	testShouldReturnErrorGiven(t,
		&Format{},
		errorMessage,
	)
}

func TestShouldReturnErrorGivenWordOfIncompatibleLength(t *testing.T) {
	const (
		errorMessage = "%[1]s error: " +
			"The length of a word should be a multiple of eight " +
			"in the range [8, 64]. " +
			"Argument to %[1]s points to a format-struct \"binary.Format\" " +
			"that has a word \"Word\" " +
			"of length 36 not in {8, 16, 24, ... 64}."
	)

	type (
		Word struct {
			BitField uint `bitfield:"36"`
		}

		Format struct {
			Word `word:"36"`
		}
	)

	testShouldReturnErrorGiven(t,
		&Format{},
		errorMessage,
	)
}

func TestShouldReturnErrorGivenWordWithNoBitFields(t *testing.T) {
	const (
		errorMessage = "%[1]s error: " +
			"A word-struct should have exported fields " +
			"representing bit fields. " +
			"Argument to %[1]s points to a format-struct \"binary.Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has no bit fields."
	)

	type (
		Word struct{}

		Format struct {
			Word `word:"32"`
		}
	)

	testShouldReturnErrorGiven(t,
		&Format{},
		errorMessage,
	)
}

func TestShouldReturnErrorGivenBitFieldOfUnsupportedType(
	t *testing.T,
) {
	const (
		errorMessage = "%[1]s error: " +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"of type uintN or bool. " +
			"Argument to %[1]s points to a format-struct \"binary.Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has a bit field \"BitField\" " +
			"of unsupported type \"int\"."
	)

	type (
		Word struct {
			BitField int `bitfield:"32"`
		}

		Format struct {
			Word `word:"32"`
		}
	)

	testShouldReturnErrorGiven(t,
		&Format{},
		errorMessage,
	)
}

func TestShouldReturnErrorGivenBitFieldWithNoStructTag(t *testing.T) {
	const (
		errorMessage = "%[1]s error: " +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"tagged with a key \"bitfield\" and a value " +
			"indicating the length of the bit field in number of bits " +
			"(e.g. `bitfield:\"1\"`). " +
			"Argument to %[1]s points to a format-struct \"binary.Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has a bit field \"BitField\" " +
			"with no struct tag."
	)

	type (
		Word struct {
			BitField uint
		}

		Format struct {
			Word `word:"32"`
		}
	)

	testShouldReturnErrorGiven(t,
		&Format{},
		errorMessage,
	)
}

func TestShouldReturnErrorGivenBitFieldWithMalformedTag(t *testing.T) {
	const (
		errorMessage = "%[1]s error: " +
			"A bit field is represented " +
			"by an exported field of a word-struct " +
			"tagged with a key \"bitfield\" and a value " +
			"indicating the length of the bit field in number of bits " +
			"(e.g. `bitfield:\"1\"`). " +
			"Argument to %[1]s points to a format-struct \"binary.Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has a bit field \"BitField\" " +
			"with a malformed struct tag."
	)

	type (
		Word struct {
			BitField uint `bitfield:"-32"`
		}

		Format struct {
			Word `word:"32"`
		}
	)

	testShouldReturnErrorGiven(t,
		&Format{},
		errorMessage,
	)
}

func TestShouldReturnErrorGivenBitFieldOfLengthOverflowingType(t *testing.T) {
	const (
		errorMessage = "%[1]s error: " +
			"The number of unique values a bit field can contain " +
			"must not exceed the size of its type. " +
			"Argument to %[1]s points to a format-struct \"binary.Format\" " +
			"nesting a word-struct \"Word\" " +
			"that has a bit field \"BitField\" " +
			"of length 32 exceeding the size of type \"uint8\"."
	)

	type (
		Word struct {
			BitField uint8 `bitfield:"32"`
		}

		Format struct {
			Word `word:"32"`
		}
	)

	testShouldReturnErrorGiven(t,
		&Format{},
		errorMessage,
	)
}

func TestShouldReturnErrorGivenWordOfLengthNotEqualToSumOfLengthsOfBitFields(
	t *testing.T,
) {
	const (
		errorMessage = "%[1]s error: " +
			"The length of a word " +
			"should be equal to the sum of lengths of its bit fields. " +
			"Argument to %[1]s points to a format-struct \"binary.Format\" " +
			"that has a word \"Word\" " +
			"of length 32 " +
			"not equal to the sum of the lengths of its bit fields, 31."
	)

	type (
		Word struct {
			BitField uint `bitfield:"31"`
		}

		Format struct {
			Word `word:"32"`
		}
	)

	testShouldReturnErrorGiven(t,
		&Format{},
		errorMessage,
	)
}

func testShouldReturnErrorGiven(t *testing.T,
	pointer interface{}, errorMessage string,
) {
	const (
		marshal   = "Marshal"
		unmarshal = "Unmarshal"
	)

	var (
		bytes []byte
		e     error
	)

	bytes, e = Marshal(pointer)

	assert.Zero(t,
		len(bytes),
	)

	assert.Equal(t,
		fmt.Sprintf(errorMessage, marshal),
		e.Error(),
	)

	e = Unmarshal(internetHeaderBytes, pointer)

	assert.Equal(t,
		fmt.Sprintf(errorMessage, unmarshal),
		e.Error(),
	)
}

func TestShouldReturnErrorGivenLengthOfByteSliceNotEqualToFormatLength(
	t *testing.T,
) {
	const (
		errorMessage = "Unmarshal error: " +
			"A byte slice into which a format-struct would be unmarshalled " +
			"should be of length equal to the sum of lengths of words " +
			"in the format represented by the struct. " +
			"Argument to Unmarshal points to a format-struct " +
			"\"binary.Format\" " +
			"of length 32 bits " +
			"not equal to the length of the byte slice, 8 bits."
	)

	type (
		Word struct {
			BitField uint `bitfield:"32"`
		}

		Format struct {
			Word `word:"32"`
		}
	)

	var (
		bytes []byte = make([]byte, 1)
		e     error
	)

	e = Unmarshal(bytes, &Format{})

	assert.Equal(t,
		errorMessage,
		e.Error(),
	)
}

const (
	destinationAddressOctet = 8
	fragmentOffset          = 8191
	headerChecksum          = 0
	identification          = 0
	sourceAddressOctet      = 1
	timeToLive              = 1
	totalLength             = 65535
)

var (
	internetHeaderStruct = rfc791.RFC791InternetHeaderFormatWithoutOptions{
		rfc791.RFC791InternetHeaderFormatWord0{
			Version:     rfc791.RFC791InternetHeaderVersion,
			IHL:         rfc791.RFC791InternetHeaderLengthWithoutOptions,
			Precedence:  rfc791.RFC791InternetHeaderPrecedenceNetworkControl,
			Delay:       rfc791.RFC791InternetHeaderDelayNormal,
			Throughput:  rfc791.RFC791InternetHeaderThroughputHigh,
			Reliability: rfc791.RFC791InternetHeaderReliabilityNormal,
			TotalLength: totalLength,
		},
		rfc791.RFC791InternetHeaderFormatWord1{
			Identification: identification,
			FlagsBit1:      rfc791.RFC791InternetHeaderFlagsBit1DoNotFragment,
			FlagsBit2:      rfc791.RFC791InternetHeaderFlagsBit2LastFragment,
			FragmentOffset: fragmentOffset,
		},
		rfc791.RFC791InternetHeaderFormatWord2{
			TimeToLive:     timeToLive,
			Protocol:       rfc791.RFC791InternetHeaderProtocolTCP,
			HeaderChecksum: headerChecksum,
		},
		rfc791.RFC791InternetHeaderFormatWord3{
			SourceAddressOctet0: sourceAddressOctet,
			SourceAddressOctet1: sourceAddressOctet,
			SourceAddressOctet2: sourceAddressOctet,
			SourceAddressOctet3: sourceAddressOctet,
		},
		rfc791.RFC791InternetHeaderFormatWord4{
			DestinationAddressOctet0: destinationAddressOctet,
			DestinationAddressOctet1: destinationAddressOctet,
			DestinationAddressOctet2: destinationAddressOctet,
			DestinationAddressOctet3: destinationAddressOctet,
		},
	}

	internetHeaderStruct1 rfc791.RFC791InternetHeaderFormatWithoutOptions

	internetHeaderStructV1p1 = v1p1.RFC791InternetHeaderFormatWithoutOptions{
		v1p1.RFC791InternetHeaderFormatWord0{
			Version:     v1p1.RFC791InternetHeaderVersion,
			IHL:         v1p1.RFC791InternetHeaderLengthWithoutOptions,
			Precedence:  v1p1.RFC791InternetHeaderPrecedenceNetworkControl,
			Delay:       v1p1.RFC791InternetHeaderDelayNormal,
			Throughput:  v1p1.RFC791InternetHeaderThroughputHigh,
			Reliability: v1p1.RFC791InternetHeaderReliabilityNormal,
			TotalLength: totalLength,
		},
		v1p1.RFC791InternetHeaderFormatWord1{
			Identification: identification,
			FlagsBit1:      v1p1.RFC791InternetHeaderFlagsBit1DoNotFragment,
			FlagsBit2:      v1p1.RFC791InternetHeaderFlagsBit2LastFragment,
			FragmentOffset: fragmentOffset,
		},
		v1p1.RFC791InternetHeaderFormatWord2{
			TimeToLive:     timeToLive,
			Protocol:       v1p1.RFC791InternetHeaderProtocolTCP,
			HeaderChecksum: headerChecksum,
		},
		v1p1.RFC791InternetHeaderFormatWord3{
			SourceAddressOctet0: sourceAddressOctet,
			SourceAddressOctet1: sourceAddressOctet,
			SourceAddressOctet2: sourceAddressOctet,
			SourceAddressOctet3: sourceAddressOctet,
		},
		v1p1.RFC791InternetHeaderFormatWord4{
			DestinationAddressOctet0: destinationAddressOctet,
			DestinationAddressOctet1: destinationAddressOctet,
			DestinationAddressOctet2: destinationAddressOctet,
			DestinationAddressOctet3: destinationAddressOctet,
		},
	}

	internetHeaderStruct1V1p1 v1p1.RFC791InternetHeaderFormatWithoutOptions

	internetHeaderBytes = []byte{
		0b01000101, 0b11101000, 0b11111111, 0b11111111,
		0b00000000, 0b00000000, 0b01011111, 0b11111111,
		0b00000001, 0b00000110, 0b00000000, 0b00000000,
		0b00000001, 0b00000001, 0b00000001, 0b00000001,
		0b00001000, 0b00001000, 0b00001000, 0b00001000,
	}
)
