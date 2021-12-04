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

This module is a drop-in replacement for `encoding/binary`
supplying the missing functions
so that developers can define custom binary formats using only struct tags and
convert between structs and byte slices
using `Marshal()` and `Unmarshal()` with their familiar signatures,
all while retaining precise, bit-level control over data structures.

## Binary Formats
Message and file formats specify how bits are arranged to encode information.
Control over individual bits or groups smaller than a byte is often required
to put together and take apart these binary structures.

### Message Formats
The following specifications
quoted from [Section 3.1](https://datatracker.ietf.org/doc/html/rfc791#section-3.1)
of RFC 791 Internet Protocol and
[Section 3.1](https://datatracker.ietf.org/doc/html/rfc793#section-3.1)
of RFC 793 Transmission Control Protocol
describe the anatomy of TCP/IP headers
at the beginning of every internet datagram ("packet").

```
 A summary of the contents of the internet header follows:


    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |Version|  IHL  |Type of Service|          Total Length         |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |         Identification        |Flags|      Fragment Offset    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  Time to Live |    Protocol   |         Header Checksum       |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                       Source Address                          |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Destination Address                        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Options                    |    Padding    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

                    Example Internet Datagram Header

                               Figure 4.

  Note that each tick mark represents one bit position.
```

It is highly unlikely that a developer
would ever need to implement these protocols,
since in Go the standard library package [`net`](https://pkg.go.dev/net)
supplies types and methods that abstract away low-level details,
but they make appropriate illustrations of binary message formats.

```
  TCP Header Format


    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |          Source Port          |       Destination Port        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                        Sequence Number                        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Acknowledgment Number                      |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  Data |           |U|A|P|R|S|F|                               |
   | Offset| Reserved  |R|C|S|S|Y|I|            Window             |
   |       |           |G|K|H|T|N|N|                               |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |           Checksum            |         Urgent Pointer        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Options                    |    Padding    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                             data                              |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

                            TCP Header Format

          Note that one tick mark represents one bit position.

                               Figure 3.
```

### File Formats
Binary file formats are not significantly different from message formats
from an application developer's perspective.
[RFC 1952](https://datatracker.ietf.org/doc/html/rfc1952)
describes the GZIP File Format Specification.

```
   1.2. Intended audience

      ...

      The text of the specification assumes a basic background in
      programming at the level of bits and other primitive data
      representations.

   2.1. Overall conventions

      In the diagrams below, a box like this:

         +---+
         |   | <-- the vertical bars might be missing
         +---+

      represents one byte; ...

   2.2. File format

      A gzip file consists of a series of "members" (compressed data
      sets).  The format of each member is specified in the following
      section.  The members simply appear one after another in the file,
      with no additional information before, between, or after them.

   2.3. Member format

      Each member has the following structure:

         +---+---+---+---+---+---+---+---+---+---+
         |ID1|ID2|CM |FLG|     MTIME     |XFL|OS | (more-->)
         +---+---+---+---+---+---+---+---+---+---+

      ...

      2.3.1. Member header and trailer

         ...

         FLG (FLaGs)
            This flag byte is divided into individual bits as follows:

               bit 0   FTEXT
               bit 1   FHCRC
               bit 2   FEXTRA
               bit 3   FNAME
               bit 4   FCOMMENT
               bit 5   reserved
               bit 6   reserved
               bit 7   reserved
```

Package [`compress/gzip`](https://pkg.go.dev/compress/gzip)
in the Go standard library spares developers the need to (de)serialise
or even to understand the GZIP file format,
but the example is included here as a stand-in for other, custom formats.

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
