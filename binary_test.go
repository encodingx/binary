package binary

import (
	"testing"

	"github.com/encodingx/binary/pkg/rfc791"
	"github.com/stretchr/testify/assert"
)

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

	internetHeaderBytes = []byte{
		0b01000101, 0b11101000, 0b11111111, 0b11111111,
		0b00000000, 0b00000000, 0b01011111, 0b11111111,
		0b00000001, 0b00000110, 0b00000000, 0b00000000,
		0b00000001, 0b00000001, 0b00000001, 0b00000001,
		0b00001000, 0b00001000, 0b00001000, 0b00001000,
	}
)

var (
	internetHeaderStruct1 rfc791.RFC791InternetHeaderFormatWithoutOptions
)

func TestMarshal(t *testing.T) {
	// Given a format-struct variable representing a binary message or file
	// And the struct field values do not overflow corresponding bit fields

	// When I pass to function Marshal() a pointer to that struct variable
	var (
		bytes []byte
		e     error
	)

	bytes, e = Marshal(&internetHeaderStruct)

	// Then Marshal() should return a slice of bytes and a nil error
	assert.Nil(t, e)

	// And I should see struct field values reflected as bits in those bytes
	assert.Equal(t,
		internetHeaderBytes, bytes,
		"Struct field values should be reflected as bits "+
			"in the slice of bytes returned by Marshal() "+
			"when it is passed a pointer to a struct variable.",
	)

	// And I should see that the lengths of the slice and the format are equal
}

func TestUnmarshal(t *testing.T) {
	// Given a format-struct type representing a binary message or file format
	// And a slice of bytes containing a binary message or file
	// And the lengths of the slice and the format (measured in bits) are equal

	// When I pass to function Unmarshal() the slice of bytes as an argument
	// And I pass to the function a pointer to the struct as a second argument
	var (
		e error
	)

	e = Unmarshal(internetHeaderBytes, &internetHeaderStruct1)

	// Then Unmarshal() should return a nil error
	assert.Nil(t, e)

	// And I should see struct field values matching the bits in those bytes
	assert.Equal(t,
		internetHeaderStruct, internetHeaderStruct1,
		"Struct field values should match the bits "+
			"in a slice of bytes passed to Unmarshal() "+
			"alongside a pointer to a struct variable.",
	)
}

/*
   Scenario:
       Given a variable that is not a pointer
       When I pass to <function> as an argument such a variable
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           Argument to <function> should be a pointer to a format-struct.
           Argument to <function> is not a pointer.
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

   Scenario:
       Given a pointer that does not point to a struct variable
       When I pass to <function> such a pointer
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           Argument to <function> should be a pointer to a format-struct.
           Argument to <function> does not point to a struct variable.
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

   Scenario:
       Given a format-struct with no exported fields
       When I pass to <function> a pointer to such a format-struct
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           A format-struct should nest exported word-structs.
           Argument to <function> points to a format-struct '[FormatStructType]'
           with no exported fields.
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

   Scenario:
       Given an exported field in a format-struct is not of type struct
       When I pass to <function> a pointer to such a format-struct
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           A format-struct should nest exported word-structs.
           Argument to <function> points to a format-struct '[FormatStructType]'
           with an exported field '[NameOfStructField]' that is not a struct.
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

   Scenario:
       Given an exported field in a format-struct with no struct tag
       When I pass to <function> a pointer to such a format-struct
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           Exported fields in a format-struct should be tagged
           with a key "word" and a value
           indicating the length of the word in number of bits
           (e.g. `word:"32"`).
           Argument to <function> points to a format-struct '[FormatStructType]'
           with an exported field '[NameOfStructField]' that has no struct tag.
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

   Scenario:
       Given an exported field in a format-struct with a malformed struct tag
           """
           A struct tag is malformed when its key is not "word"
           or when its value cannot be parsed as an unsigned integer.
           """
       When I pass to <function> a pointer to such a format-struct
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           Exported fields in a format-struct should be tagged
           with a key "word" and a value
           indicating the length of the word in number of bits
           (e.g. `word:"32"`).
           Argument to <function> points to a format-struct '[FormatStructType]'
           with an exported field '[NameOfStructField]'
           that has a malformed struct tag: [message of wrapped error].
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

   Scenario:
       Given a word of length not a multiple of eight in the range [8, 64]
       When I pass to <function> a pointer to a format-struct nesting such word
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           The length of a word should be a multiple of eight
           in the range [8, 64].
           Argument to <function> points to a format-struct '[FormatStructType]'
           containing a word-struct '[NameOfStructField]' of length [length]
           not in {8, 16, 24, ... 64}.
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

   Scenario:
       Given a word-struct containing no exported fields
       When I pass to <function> a pointer to a format-struct nesting such word
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           A word-struct should nest exported fields representing bit fields.
           Argument to <function> points to a format-struct '[FormatStructType]'
           containing a word-struct '[FieldName]' of type '[WordStructType]',
           which has no exported fields.
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

   Scenario:
       Given a word-struct containing a field of unsupported type
           """
           Supported types are uint, uintN where N = {8, 16, 32, 64} and bool.
           """
       When I pass to <function> a pointer to a format-struct nesting such word
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           The fields of a word-struct should be of type uintN or bool.
           Argument to <function> points to a format-struct '[FormatStructType]'
           containing a word-struct '[FieldName0]' of type '[WordStructType]',
           which has a field '[FieldName1]' of unsupported type '[FieldType]'.
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

   Scenario:
       Given a word-struct containing a field with no struct tag
       When I pass to <function> a pointer to a format-struct nesting such word
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           Exported fields in a word-struct should be tagged
           with a key "bitfield" and a value
           indicating the length of the bit field in number of bits
           (e.g. `bitfield:"1"`).
           Argument to <function> points to a format-struct '[FormatStructType]'
           containing a word-struct '[FieldName0]' of type '[WordStructType]',
           which has a field '[FieldName1]' that has no struct tag.
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

   Scenario:
       Given a word-struct containing a field with a malformed struct tag
           """
           A struct tag is malformed when its key is not "bitfield"
           or when its value cannot be parsed as an unsigned integer.
           """
       When I pass to <function> a pointer to a format-struct nesting such word
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           Exported fields in a word-struct should be tagged
           with a key "bitfield" and a value
           indicating the length of the bit field in number of bits
           (e.g. `bitfield:"1"`).
           Argument to <function> points to a format-struct '[FormatStructType]'
           containing a word-struct '[FieldName0]' of type '[WordStructType]',
           which has a field '[FieldName1]' that has a malformed struct tag:
           [message of wrapped error].
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

   Scenario:
       Given a word-struct containing a field with a v1.1-style struct tag
           """
           type RFC791InternetHeaderFormatWord0 struct {
               Version     uint8  `bitfield:"4,28"`
               IHL         uint8  `bitfield:"4,24"`
               Precedence  uint8  `bitfield:"3,21"`
               Delay       bool   `bitfield:"1,20"`
               Throughput  bool   `bitfield:"1,19"`
               Reliability bool   `bitfield:"1,18"`
               Reserved    uint8  `bitfield:"2,16"`
               TotalLength uint16 `bitfield:"16,0"`
           }
           """
       When I pass to <function> a pointer to a format-struct nesting such word
       Then <function> should return <a byte slice and> a nil error

       Examples:
           | function  | a byte slice and |
           | Marshal   | a byte slice and |
           | Unmarshal |                  |

   Scenario:
       Given a word-struct with a bit field of length overflowing its type
           """
           A bit field overflows a type
           when it is long enough to represent values
           outside the set of values of the type.
           """
       When I pass to <function> a pointer to a format-struct nesting such word
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           The number of unique values a bit field can contain
           must not exceed the size of its type.
           Argument to <function> points to a format-struct '[FormatStructType]'
           containing a word-struct '[FieldName0]' of type '[WordStructType]',
           which has a bit field '[FieldName1]' of length [length]
           exceeding the size of type [FieldType], [size].
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

   Scenario:
       Given a word of length not equal to the sum of lengths of its bit fields
       When I pass to <function> a pointer to a format-struct nesting such word
       Then <function> should return <a byte slice of zero length and> an error
           """
           <function> error:
           The length of a word
           should be equal to the sum of lengths of its bit fields.
           Argument to <function> points to a format-struct '[NameOfStructType]'
           containing a word-struct '[NameOfStructField]' of length [length]
           not equal to the sum of the lengths of its bit fields, [sum].
           """

       Examples:
           | function  | a byte slice of zero length and |
           | Marshal   | a byte slice of zero length and |
           | Unmarshal |                                 |

    Scenario:
        Given a format-struct type representing a binary message or file format
        And a slice of bytes containing a binary message or file
        And the lengths of the slice and the format (in bits) are not equal
        When I pass to function Unmarshal() the slice of bytes as an argument
        And I pass to the function a pointer to the struct as a second argument
        Then Unmarshal() should return an error
            """
            Unmarshal error:
            A byte slice into which a format-struct would be unmarshalled
            should have the same length as the format represented by the struct.
            The length of a format is the sum of the lengths of the words in it.
            Argument to Unmarshal() points to a format-struct '[NameOfStructType]'
            of length [length0] not equal to length [length1] of the byte slice.
            (Lengths are measured in number of bits.)
            """
*/

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
