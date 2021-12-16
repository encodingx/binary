# Marshal and Unmarshal Binary Formats in Go
The Go standard library features packages for marshalling structs
into text-based message and file formats, and vice versa.
Packages `encoding/json` and `encoding/xml` are well known
for the convenient functions `Marshal()` and `Unmarshal()` they provide
that leverage the declarative nature of struct tags.
These functions are however missing from
[`encoding/binary`](https://pkg.go.dev/encoding/binary),
leaving developers without an accessible and intuitive way
to work with binary formats.

This module supplies the missing functions
so that developers can define custom binary formats using only struct tags and
convert between structs and byte slices
using `Marshal()` and `Unmarshal()` with their familiar signatures.

## Binary Formats
Message and file formats specify how bits are arranged to encode information.
Control over individual bits or groups smaller than a byte is often required
to put together and take apart these binary structures.

### Message Formats
Describing the anatomy of TCP/IP headers
at the beginning of every internet datagram ("packet")
are
* [Section 3.1](https://datatracker.ietf.org/doc/html/rfc791#section-3.1)
of RFC 791 Internet Protocol, and
* [Section 3.1](https://datatracker.ietf.org/doc/html/rfc793#section-3.1)
of RFC 793 Transmission Control Protocol

### File Formats
Binary file formats are not significantly different from message formats
from an application developer's perspective.
[RFC 1952](https://datatracker.ietf.org/doc/html/rfc1952)
describes the GZIP File Format Specification.

## Working with Binary Formats in Go
The smallest data structures Go provides are
the basic type `byte` (alias of `uint8`, an unsigned 8-bit integer), and `bool`,
both eight bits long.
To manipulate data at a scale smaller than eight bits
would require the use of bitwise logical and shift [operators](https://go.dev/ref/spec#Arithmetic_operators) such as
[AND](https://en.wikipedia.org/wiki/Bitwise_operation#AND) (`&`),
[OR](https://en.wikipedia.org/wiki/Bitwise_operation#OR) (`|`),
left [shift](https://en.wikipedia.org/wiki/Bitwise_operation#Bit_shifts) (`<<`),
and right shift (`>>`).

### Relevant Questions Posted on StackOverflow
Suggestions on StackOverflow are limited to the use of bitwise operators.

* [Golang: Parse bit values from a byte](https://stackoverflow.com/questions/54809254/golang-parse-bit-values-from-a-byte)
* [Creating 8 bit binary data from 4,3, and 1 bit data in Golang](https://stackoverflow.com/questions/61885659/creating-8-bit-binary-data-from-4-3-and-1-bit-data-in-golang)
* [How to pack the C bit field struct via encoding package in GO?](https://stackoverflow.com/questions/60180827/how-to-pack-the-c-bit-field-struct-via-encoding-package-in-go)

## Behaviour-Driven Specifications
The following is an excerpt from the [specifications](docs/binary.feature)
of this module.

```gherkin
Feature: Marshal and Unmarshal

    As a Go developer implementing a binary message or file format,
    I want a pair of functions "Marshal/Unmarshal" like those in "encoding/json"
    that convert a struct into a series of bits in a byte slice and vice versa,
    so that I can avoid the complexities of custom bit manipulation.

    Background:
        # Ubiquitous language
        Given a message or file "format"
            """
            A format specifies how bits are arranged to encode information.
            """
        And the format is a series of "bit fields"
            """
            A bit field is one or more adjacent bits representing a value,
            and should not be confused with struct fields.
            """
        And adjacent bit fields are grouped into "words"
            """
            A word is a series of bits that can be simultaneously processed
            by a given computer architecture and programming language.
            """

        # Define format-structs
        And a format is represented by a type definition of a "format-struct"
        And the format-struct nests one or more exported "word-structs"
        And the words are tagged to indicate their lengths in number of bits
            """
            type RFC791InternetHeaderFormatWithoutOptions struct {
                RFC791InternetHeaderFormatWord0 `word:"32"`
                RFC791InternetHeaderFormatWord1 `word:"32"`
                RFC791InternetHeaderFormatWord2 `word:"32"`
                RFC791InternetHeaderFormatWord3 `word:"32"`
                RFC791InternetHeaderFormatWord4 `word:"32"`
            }
            """
        And the length of each word is a multiple of eight in the range [8, 64]

        # Define word-structs
        And each word-struct has exported field(s) corresponding to bit field(s)
        And the fields are of unsigned integer or boolean types
        And the fields are tagged to indicate the lengths of those bit fields
            """
            type RFC791InternetHeaderFormatWord0 struct {
                Version     uint8  `bitfield:"4"`
                IHL         uint8  `bitfield:"4"`
                Precedence  uint8  `bitfield:"3"`
                Delay       bool   `bitfield:"1"`
                Throughput  bool   `bitfield:"1"`
                Reliability bool   `bitfield:"1"`
                Reserved    uint8  `bitfield:"2"`
                TotalLength uint16 `bitfield:"16"`
            }
            """
        And the length of each bit field does not overflow the type of the field
            """
            A bit field overflows a type
            when it is long enough to represent values
            outside the set of values of the type.
            """
        And the sum of lengths of all fields is equal to the length of that word

    Scenario: Marshal a struct into a byte slice
        Given a format-struct variable representing a binary message or file
            """
            internetHeader = RFC791InternetHeaderFormatWithoutOptions{
                RFC791InternetHeaderFormatWord0{
                    Version: 4,
                    IHL:     5,
                    // ...
                },
                // ...
            }
            """
        And the struct field values do not overflow corresponding bit fields
            """
            A struct field value overflows its corresponding bit field
            when it falls outside the range of values
            that can be represented by that bit field given its length.
            """
        When I pass to function Marshal() a pointer to that struct variable
            """
            var (
                bytes []byte
                e     error
            )

            bytes, e = binary.Marshal(&internetHeader)
            """
        Then Marshal() should return a slice of bytes and a nil error
        And I should see struct field values reflected as bits in those bytes
            """
            log.Printf("%08b", bytes)
            // [01000101 ...]

            log.Println(e == nil)
            // true
            """

    Scenario: Unmarshal a byte slice into a struct
        Given a format-struct type representing a binary message or file format
            """
            var internetHeader RFC791InternetHeaderFormatWithoutOptions
            """
        And a slice of bytes containing a binary message or file
            """
            var bytes []byte

            // ...

            log.Printf("%08b", bytes)
            // [01000101 ...]
            """
        And the lengths of the slice and the format (measured in bits) are equal
            """
            The length of a format is the sum of lengths of the words in it.
            The length of a word is the sum of lengths of the bit fields in it.
            """
        When I pass to function Unmarshal() the slice of bytes as an argument
        And I pass to the function a pointer to the struct as a second argument
            """
            e = binary.Unmarshal(bytes, &internetHeader)
            """
        Then Unmarshal() should return a nil error
        And I should see struct field values matching the bits in those bytes
            """
            log.Println(e == nil)
            // true

            log.Println(internetHeader.RFC791InternetHeaderFormatWord0.Version)
            // 4

            log.Println(internetHeader.RFC791InternetHeaderFormatWord0.IHL)
            // 5
            """
```

See the rest of the [specifications](docs/binary.feature) for error scenarios.

## Using This Module
Encoding and decoding binary formats in Go should be a matter of
attaching tags to struct fields and calling `Marshal()`/`Unmarshal()`,
in the same fashion as JSON (un)marshalling familiar to many Go developers.
This module is meant to be imported in lieu of `encoding/binary`.
Exported variables, types and functions of the standard library package
pass through and are available as though that package had been imported.

### Structs
#### Words
A word is made of one or more fields of up to 64 bits in length.
The length and offset (see [definition](#offset)) of a field in number of bits
must be indicated by a struct tag of the format `bitfield:"<length>,<offset>"`.
There should be no gaps nor overlaps between fields, and
the sum of the lengths of all fields in a word
must be equal to the word length declared in the [format struct](#formats).
Unused or "reserved" fields should nonetheless by defined and tagged
even though they contain all zeroes.

```go
package demo

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
```

As of Version 1, supported field types are
`uint`, `uint8` a.k.a. `byte`, `uint16`, `uint32`, `uint64` and `bool`.
[Defined types](https://go.dev/ref/spec#Type_definitions)
(e.g. `type RFC791InternetHeaderPrecedence uint8`)
having the above underlying types are compatible.

A word struct cannot be marshalled/unmarshalled by itself;
it must be a field in a format struct, explained below.

##### Offset
`<offset>` is an integer
representing the number of places the bit field should be shifted left
from the rightmost section of a word for its position to be appropriate.
It is also the number of places to the right of the rightmost bit of the field.
The offset of every field is the sum of the lengths of all fields to its right.

#### Formats
A "format" is a struct that represents a binary message or file format,
made up of one or more "words" (see [section](#words) on words below).
A format struct must have one or more fields,
all of which must be structs and bear a tag of the format `word:"<length>"`,
where `<length>` is an integer multiple of eight (up to a limit of 64)
indicating the number of bits in the word.

```go
package demo

type RFC791InternetHeaderFormatWithoutOptions struct {
	RFC791InternetHeaderFormatWord0 `word:"32"`
	RFC791InternetHeaderFormatWord1 `word:"32"`
	RFC791InternetHeaderFormatWord2 `word:"32"`
	RFC791InternetHeaderFormatWord3 `word:"32"`
	RFC791InternetHeaderFormatWord4 `word:"32"`
}
```

`Marshal()` and `Unmarshal()` expect a pointer to a format struct
as an argument.

### Example
```go
package main

import (
	"log"

	"github.com/encodingx/binary"
	"github.com/encodingx/binary/pkg/demo"
)

func main() {
	const (
		destinationAddressOctet = 8
		fragmentOffset          = 8191
		headerChecksum          = 0
		identification          = 0
		sourceAddressOctet      = 1
		timeToLive              = 1
		totalLength             = 65535

		formatBytes = "%08b"
	)

	var (
		struct0 demo.RFC791InternetHeaderFormatWithoutOptions
		struct1 demo.RFC791InternetHeaderFormatWithoutOptions

		bytes []byte
		e     error
	)

	struct0 = demo.RFC791InternetHeaderFormatWithoutOptions{
		demo.RFC791InternetHeaderFormatWord0{
			Version:     demo.RFC791InternetHeaderVersion,
			IHL:         demo.RFC791InternetHeaderLengthWithoutOptions,
			Precedence:  demo.RFC791InternetHeaderPrecedenceNetworkControl,
			Delay:       demo.RFC791InternetHeaderDelayNormal,
			Throughput:  demo.RFC791InternetHeaderThroughputHigh,
			Reliability: demo.RFC791InternetHeaderReliabilityNormal,
			TotalLength: totalLength,
		},
		demo.RFC791InternetHeaderFormatWord1{
			Identification: identification,
			FlagsBit1:      demo.RFC791InternetHeaderFlagsBit1DoNotFragment,
			FlagsBit2:      demo.RFC791InternetHeaderFlagsBit2LastFragment,
			FragmentOffset: fragmentOffset,
		},
		demo.RFC791InternetHeaderFormatWord2{
			TimeToLive:     timeToLive,
			Protocol:       demo.RFC791InternetHeaderProtocolTCP,
			HeaderChecksum: headerChecksum,
		},
		demo.RFC791InternetHeaderFormatWord3{
			SourceAddressOctet0: sourceAddressOctet,
			SourceAddressOctet1: sourceAddressOctet,
			SourceAddressOctet2: sourceAddressOctet,
			SourceAddressOctet3: sourceAddressOctet,
		},
		demo.RFC791InternetHeaderFormatWord4{
			DestinationAddressOctet0: destinationAddressOctet,
			DestinationAddressOctet1: destinationAddressOctet,
			DestinationAddressOctet2: destinationAddressOctet,
			DestinationAddressOctet3: destinationAddressOctet,
		},
	}

	bytes, e = binary.Marshal(&struct0)
	if e != nil {
		log.Fatalln(e)
	}

	log.Printf(formatBytes, bytes)
	// [01000101 11101000 11111111 11111111
	//  00000000 00000000 01011111 11111111
	//  00000001 00000110 00000000 00000000
	//  00000001 00000001 00000001 00000001
	//  00001000 00001000 00001000 00001000]

	e = binary.Unmarshal(bytes, &struct1)
	if e != nil {
		log.Fatalln(e)
	}

	log.Println(struct1 == struct0)
	// true
}
```

Compare the first word (four bytes) in the output of `Marshal()`
to the specifications for the first 32 bits of the Internet Header.

```
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |Version|  IHL  |Type of Service|          Total Length         |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```
```
  Version:  4 bits

    The Version field indicates the format of the internet header.  This
    document describes version 4.

  IHL:  4 bits

    Internet Header Length is the length of the internet header in 32
    bit words, and thus points to the beginning of the data.  Note that
    the minimum value for a correct header is 5.

  Type of Service:  8 bits

    The Type of Service provides an indication of the abstract
    parameters of the quality of service desired.  These parameters are
    to be used to guide the selection of the actual service parameters
    when transmitting a datagram through a particular network.  Several
    networks offer service precedence, which somehow treats high
    precedence traffic as more important than other traffic (generally
    by accepting only traffic above a certain precedence at time of high
    load).  The major choice is a three way tradeoff between low-delay,
    high-reliability, and high-throughput.

      Bits 0-2:  Precedence.
      Bit    3:  0 = Normal Delay,       1 = Low Delay.
      Bits   4:  0 = Normal Throughput,  1 = High Throughput.
      Bits   5:  0 = Normal Reliability, 1 = High Reliability.
      Bit  6-7:  Reserved for Future Use.

         0     1     2     3     4     5     6     7
      +-----+-----+-----+-----+-----+-----+-----+-----+
      |                 |     |     |     |     |     |
      |   PRECEDENCE    |  D  |  T  |  R  |  0  |  0  |
      |                 |     |     |     |     |     |
      +-----+-----+-----+-----+-----+-----+-----+-----+

        Precedence

          111 - Network Control
          110 - Internetwork Control
          101 - CRITIC/ECP
          100 - Flash Override
          011 - Flash
          010 - Immediate
          001 - Priority
          000 - Routine
```

The rest of the document is omitted for brevity.
Similar excerpts from Section 3.1 of RFC 791 are quoted in comments to
the definition of struct `demo.RFC791InternetHeaderFormatWithoutOptions`.
Values of constants in the struct literal in the example code above
are declared in the same file containing the struct definition.

## Performance and Optimisation
This module has been optimised for performance.
Suggestions to improve are welcome.

```bash
$ go test -cpuprofile cpu.prof -memprofile mem.prof -bench . -benchmem
```
```
goos: linux
goarch: amd64
pkg: github.com/encodingx/binary
cpu: Intel(R) Core(TM) i5-8259U CPU @ 2.30GHz
BenchmarkMarshal   	 1206807	       851.9 ns/op	      64 B/op	       6 allocs/op
BenchmarkUnmarshal 	 1376592	       876.0 ns/op	      64 B/op	       8 allocs/op
PASS
ok  	github.com/encodingx/binary	4.015s
```
