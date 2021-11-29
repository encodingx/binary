# Encode and Decode Binary Formats in Go
This module wraps the package
[`encoding/binary`](https://pkg.go.dev/encoding/binary)
of the Go standard library and
provides the missing `Marshal()` and `Unmarshal()` functions
a la `encoding/json` and `encoding/xml`.

## Intended Audience
This module is useful to developers
implementing binary message and file format specifications
using the Go programming language.
Whereas stable implementations of most open-source formats
are readily available,
proprietary formats often require bespoke solutions.
The Go standard library provides convenient functions
`Marshal()` and `Unmarshal()`
for converting Go structs into text-based data formats
(such as [JSON](https://pkg.go.dev/encoding/json#Marshal) and
[XML](https://pkg.go.dev/encoding/xml#Marshal))
and vice versa, but
their counterparts are sorely missing from the `encoding/binary`.

## Binary Message and File Formats
Message and file formats specify how bits are arranged to encode information.
Control over individual bits or groups smaller than a byte is often required
to put together and take apart these binary structures.

### Message Formats
Messages are the lifeblood of network applications.
The following specifications describe the anatomy of TCP/IP headers
at the beginning of every internet datagram (a.k.a. "packet").

#### RFC 791 Internet Protocol
Section [3.1](https://datatracker.ietf.org/doc/html/rfc791#section-3.1)
Internet Header Format
```
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

  ...

  Type of Service:  8 bits

    ...

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

Notice how the Internet Datagram Header is visualised as a series
of five or more 32-bit "words".

#### RFC 793 Transmission Control Protocol
Section [3.1](https://datatracker.ietf.org/doc/html/rfc793#section-3.1)
Header Format
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

#### RFC 1952 GZIP File Format Specification
Example taken from [RFC 1952](https://datatracker.ietf.org/doc/html/rfc1952#page-5)
defining the GZIP file format:

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

## Working with Binary Formats in Go
The smallest data structures Go provides are
the basic type `byte` (alias of `uint8`), and `bool`.
Both are eight bits long.
To manipulate data at the level of bits
would require the use of bitwise logical and shift [operators](https://go.dev/ref/spec#Arithmetic_operators)
`&`, `|`, `^`, `&^`, `<<`, and `>>`.

### Relevant Questions Posted on StackOverflow
Suggestions on StackOverflow in response to the above-described need
are limited to the use of bitwise operators.

* [Golang: Parse bit values from a byte](https://stackoverflow.com/questions/54809254/golang-parse-bit-values-from-a-byte)
* [Creating 8 bit binary data from 4,3, and 1 bit data in Golang](https://stackoverflow.com/questions/61885659/creating-8-bit-binary-data-from-4-3-and-1-bit-data-in-golang)

## An Elegant Solution
Encoding and decoding binary formats in Go should be a matter of
attaching tags to struct fields and calling `Marshal()`/`Unmarshal()`.

```go
package main

import (
	"log"

	"github.com/joel-ling/go-bitfields/pkg/encoding/binary"
)

type rfc791InternetHeaderWord0 struct {
	// Reference: RFC 791 Internet Protocol, Section 3.1 Internet Header Format
	// https://datatracker.ietf.org/doc/html/rfc791#section-3.1

	Version                  uint                           `bitfield:"4,28"`
	InternetHeaderLength     uint                           `bitfield:"4,24"`
	TypeOfServicePrecedence  rfc791InternetHeaderPrecedence `bitfield:"3,21"`
	TypeOfServiceDelay       bool                           `bitfield:"1,20"`
	TypeOfServiceThroughput  bool                           `bitfield:"1,19"`
	TypeOfServiceReliability bool                           `bitfield:"1,18"`
	TypeOfServiceReserved    uint                           `bitfield:"2,16"`
	TotalLength              uint                           `bitfield:"16,0"`
}

type rfc791InternetHeaderPrecedence byte

const (
	rfc791InternetHeaderPrecedenceRoutine rfc791InternetHeaderPrecedence = iota
	rfc791InternetHeaderPrecedencePriority
	rfc791InternetHeaderPrecedenceImmediate
	rfc791InternetHeaderPrecedenceFlash
	rfc791InternetHeaderPrecedenceFlashOverride
	rfc791InternetHeaderPrecedenceCRITICOrECP
	rfc791InternetHeaderPrecedenceInternetworkControl
	rfc791InternetHeaderPrecedenceNetworkControl
)

func main() {
	const (
		internetHeaderLength  = 5
		internetHeaderVersion = 4

		internetHeaderDelayNormal = false
		internetHeaderDelayLow    = true

		internetHeaderThroughputNormal = false
		internetHeaderThroughputHigh   = true

		internetHeaderReliabilityNormal = false
		internetHeaderReliabilityHigh   = true

		internetHeaderTotalLength = 30222

		reserved = 0
	)

	const (
		formatBytes  = "%08b"
		formatStruct = "%#v"
	)

	var (
		structure0 = rfc791InternetHeaderWord0{
			Version:                  internetHeaderVersion,
			InternetHeaderLength:     internetHeaderLength,
			TypeOfServicePrecedence:  rfc791InternetHeaderPrecedenceImmediate,
			TypeOfServiceDelay:       internetHeaderDelayLow,
			TypeOfServiceThroughput:  internetHeaderThroughputNormal,
			TypeOfServiceReliability: internetHeaderReliabilityHigh,
			TypeOfServiceReserved:    reserved,
			TotalLength:              internetHeaderTotalLength,
		}

		structure1 = rfc791InternetHeaderWord0{}
	)

	var (
		bytes []byte
		e     error
	)

	bytes, e = binary.Marshal(&structure0)
	if e != nil {
		log.Fatalln(e)
	}

	log.Printf(formatBytes, bytes)
	// [01000101 01010100 01110110 00001110]

	e = binary.Unmarshal(bytes, &structure1)
	if e != nil {
		log.Fatalln(e)
	}

	log.Printf(formatStruct, structure1 == structure0)
	// true
}
```

Large structures can be encoded and decoded in 32-bit instalments.

### Struct Tags
The encoding of every struct field is determined by a struct tag
of the following key and format:

`bitfield:"<length>,<offset>"`

#### Length
`<length>` refers to an integer representing the bit-length of a field, i.e.
the number of bits in the field.

#### Offset
`<offset>` refers to an integer
representing the number of places the bit field should be shifted left
from the rightmost section of a 32-bit sequence
for its position in that sequence to be appropriate.
It is also the number of places to the right of the rightmost bit of the field.

Notice how in a fully utilised 32-bit word,
the offset of every field is the cumulative sum of their lengths,
when viewed from the rightmost field to the leftmost.
