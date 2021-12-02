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

As of Version 1.0.0, supported field types are
`uint`, `uint8` a.k.a. `byte`, `uint16`, `uint32`, `uint64` and `bool`.

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
(pprof) top 15 -cum
(pprof) top 15 -cum
Showing nodes accounting for 1.68s, 43.75% of 3.84s total
Dropped 30 nodes (cum <= 0.02s)
Showing top 15 nodes out of 70
      flat  flat%   sum%        cum   cum%
         0     0%     0%      3.75s 97.66%  testing.(*B).launch
         0     0%     0%      3.75s 97.66%  testing.(*B).runN
     0.04s  1.04%  1.04%         2s 52.08%  github.com/joel-ling/go-bitfields/pkg/codecs.(*v1Codec).Unmarshal
         0     0%  1.04%         2s 52.08%  github.com/joel-ling/go-bitfields/pkg/encoding/binary.BenchmarkUnmarshal
         0     0%  1.04%         2s 52.08%  github.com/joel-ling/go-bitfields/pkg/encoding/binary.Unmarshal (inline)
     0.22s  5.73%  6.77%      1.89s 49.22%  github.com/joel-ling/go-bitfields/pkg/codecs.unmarshalFormat
     0.02s  0.52%  7.29%      1.75s 45.57%  github.com/joel-ling/go-bitfields/pkg/codecs.(*v1Codec).Marshal
         0     0%  7.29%      1.75s 45.57%  github.com/joel-ling/go-bitfields/pkg/encoding/binary.BenchmarkMarshal
         0     0%  7.29%      1.75s 45.57%  github.com/joel-ling/go-bitfields/pkg/encoding/binary.Marshal (inline)
     0.10s  2.60%  9.90%      1.63s 42.45%  github.com/joel-ling/go-bitfields/pkg/codecs.marshalFormat
     0.29s  7.55% 17.45%      1.59s 41.41%  github.com/joel-ling/go-bitfields/pkg/codecs.unmarshalWord
     0.18s  4.69% 22.14%      1.38s 35.94%  github.com/joel-ling/go-bitfields/pkg/codecs.marshalWord
     0.05s  1.30% 23.44%      0.73s 19.01%  runtime.makeslice
     0.36s  9.38% 32.81%      0.68s 17.71%  runtime.mallocgc
     0.42s 10.94% 43.75%      0.59s 15.36%  github.com/joel-ling/go-bitfields/pkg/codecs.unmarshalField
```

Profiling reveals that about 20% of CPU time is spent on slice allocation.
