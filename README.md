# Marshal and Unmarshal Binary Formats in Go
This module is a drop-in replacement for the standard library package
[`encoding/binary`](https://pkg.go.dev/encoding/binary)
of the Go programming language and
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
and vice versa,
but their counterparts for binary formats
are sorely missing from the package `encoding/binary`.

## Working with Binary Formats in Go
Message and file formats specify how bits are arranged to encode information.
Control over individual bits or groups smaller than a byte is often required
to put together and take apart these binary structures.

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

## Using This Module
Encoding and decoding binary formats in Go should be a matter of
attaching tags to struct fields and calling `Marshal()`/`Unmarshal()`,
in the same fashion as JSON (un)marshalling familiar to many Go developers.
This module is meant to be imported in lieu of `encoding/binary`.
Exported variables, types and functions of the standard library package
pass through and are available as though that package had been imported.

```go
package main

import (
	"log"

	"github.com/joel-ling/go-bitfields/pkg/encoding/binary"
	"github.com/joel-ling/go-bitfields/pkg/structs/demo"
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

		formatBytes  = "%08b"
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
to the specifications for bit fields in that word
in [Section 3.1](https://datatracker.ietf.org/doc/html/rfc791#section-3.1)
"Internet Header Format"
of RFC 791 defining the Internet Protocol.

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
      Bit    3:  0 = Normal Delay,      1 = Low Delay.
      Bits   4:  0 = Normal Throughput, 1 = High Throughput.
      Bits   5:  0 = Normal Relibility, 1 = High Relibility.
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
the [definition](https://github.com/joel-ling/go-bitfields/blob/v1.0.2/internal/structs/demo/rfc-791-internet-header-format.go)
of struct `demo.RFC791InternetHeaderFormatWithoutOptions`.
Values of constants in the struct literal in the example code above
are declared in the same file containing the struct definition.

### Structs
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

##### Offset
`<offset>` is an integer
representing the number of places the bit field should be shifted left
from the rightmost section of a word for its position to be appropriate.
It is also the number of places to the right of the rightmost bit of the field.
The offset of every field is the sum of the lengths of all fields to its right.

### Performance
```bash
pkg/encoding/binary$ go test -cpuprofile cpu.prof -memprofile mem.prof -bench . -benchmem
```
```
goos: linux
goarch: amd64
pkg: github.com/joel-ling/go-bitfields/pkg/encoding/binary
cpu: Intel(R) Core(TM) i5-8259U CPU @ 2.30GHz
BenchmarkMarshal   	 1211152	       845.1 ns/op	      64 B/op	       6 allocs/op
BenchmarkUnmarshal 	 1370266	       878.1 ns/op	      64 B/op	       8 allocs/op
PASS
ok  	github.com/joel-ling/go-bitfields/pkg/encoding/binary	4.200s
```

#### CPU Profiling
```bash
pkg/encoding/binary$ go tool pprof cpu.prof
```
```
(pprof) top 20 -cum
Showing nodes accounting for 2.56s, 68.09% of 3.76s total
Dropped 36 nodes (cum <= 0.02s)
Showing top 20 nodes out of 75
      flat  flat%   sum%        cum   cum%
         0     0%     0%      3.65s 97.07%  testing.(*B).launch
         0     0%     0%      3.65s 97.07%  testing.(*B).runN
     0.01s  0.27%  0.27%      2.02s 53.72%  github.com/joel-ling/go-bitfields/pkg/encoding/binary.BenchmarkUnmarshal
     0.01s  0.27%  0.53%      2.01s 53.46%  github.com/joel-ling/go-bitfields/internal/codecs.(*v1Codec).Unmarshal
         0     0%  0.53%      2.01s 53.46%  github.com/joel-ling/go-bitfields/pkg/encoding/binary.Unmarshal (inline)
     0.11s  2.93%  3.46%      1.97s 52.39%  github.com/joel-ling/go-bitfields/internal/codecs.unmarshalFormat
     0.31s  8.24% 11.70%      1.78s 47.34%  github.com/joel-ling/go-bitfields/internal/codecs.unmarshalWord
     0.02s  0.53% 12.23%      1.62s 43.09%  github.com/joel-ling/go-bitfields/internal/codecs.(*v1Codec).Marshal
         0     0% 12.23%      1.62s 43.09%  github.com/joel-ling/go-bitfields/pkg/encoding/binary.BenchmarkMarshal
         0     0% 12.23%      1.62s 43.09%  github.com/joel-ling/go-bitfields/pkg/encoding/binary.Marshal (inline)
     0.15s  3.99% 16.22%      1.57s 41.76%  github.com/joel-ling/go-bitfields/internal/codecs.marshalFormat
     0.23s  6.12% 22.34%      1.16s 30.85%  github.com/joel-ling/go-bitfields/internal/codecs.marshalWord
     0.11s  2.93% 25.27%      0.91s 24.20%  runtime.makeslice
     0.37s  9.84% 35.11%      0.80s 21.28%  runtime.mallocgc
     0.30s  7.98% 43.09%      0.60s 15.96%  github.com/joel-ling/go-bitfields/internal/codecs.unmarshalField
     0.32s  8.51% 51.60%      0.42s 11.17%  github.com/joel-ling/go-bitfields/internal/codecs.marshalField
     0.28s  7.45% 59.04%      0.40s 10.64%  reflect.Value.Field
     0.19s  5.05% 64.10%      0.19s  5.05%  github.com/joel-ling/go-bitfields/internal/structs/words.(*word).Field
     0.15s  3.99% 68.09%      0.18s  4.79%  reflect.Value.SetUint
         0     0% 68.09%      0.11s  2.93%  runtime.(*mcache).nextFree
```

Profiling reveals that about 20% of CPU time is spent on slice allocation.
Accessing struct field values consumes about 15%.
