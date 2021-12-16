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
[RFC 1952](https://datatracker.ietf.org/doc/html/rfc1952#section-2)
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
```
```go
            type RFC791InternetHeaderFormatWithoutOptions struct {
                RFC791InternetHeaderFormatWord0 `word:"32"`
                RFC791InternetHeaderFormatWord1 `word:"32"`
                RFC791InternetHeaderFormatWord2 `word:"32"`
                RFC791InternetHeaderFormatWord3 `word:"32"`
                RFC791InternetHeaderFormatWord4 `word:"32"`
            }
```
```gherkin
        And the length of each word is a multiple of eight in the range [8, 64]

        # Define word-structs
        And each word-struct has exported field(s) corresponding to bit field(s)
        And the fields are of unsigned integer or boolean types
        And the fields are tagged to indicate the lengths of those bit fields
```
```go
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
```
```gherkin
        And the length of each bit field does not overflow the type of the field
            """
            A bit field overflows a type
            when it is long enough to represent values
            outside the set of values of the type.
            """
        And the sum of lengths of all fields is equal to the length of that word

    Scenario: Marshal a struct into a byte slice
        Given a format-struct variable representing a binary message or file
```
```go
            internetHeader = RFC791InternetHeaderFormatWithoutOptions{
                RFC791InternetHeaderFormatWord0{
                    Version: 4,
                    IHL:     5,
                    // ...
                },
                // ...
            }
```
```gherkin
        And the struct field values do not overflow corresponding bit fields
            """
            A struct field value overflows its corresponding bit field
            when it falls outside the range of values
            that can be represented by that bit field given its length.
            """
        When I pass to function Marshal() a pointer to that struct variable
```
```go
            var (
                bytes []byte
                e     error
            )

            bytes, e = binary.Marshal(&internetHeader)
```
```gherkin
        Then Marshal() should return a slice of bytes and a nil error
        And I should see struct field values reflected as bits in those bytes
```
```go
            log.Printf("%08b", bytes)
            // [01000101 ...]

            log.Println(e == nil)
            // true
```
```gherkin
        And I should see that the lengths of the slice and the format are equal
            """
            The length of a format is the sum of lengths of the words in it.
            The length of a word is the sum of lengths of the bit fields in it.
            """

    Scenario: Unmarshal a byte slice into a struct
        Given a format-struct type representing a binary message or file format
```
```go
            var internetHeader RFC791InternetHeaderFormatWithoutOptions
```
```gherkin
        And a slice of bytes containing a binary message or file
```
```go
            var bytes []byte

            // ...

            log.Printf("%08b", bytes)
            // [01000101 ...]
```
```gherkin
        And the lengths of the slice and the format (measured in bits) are equal
        When I pass to function Unmarshal() the slice of bytes as an argument
        And I pass to the function a pointer to the struct as a second argument
```
```go
            e = binary.Unmarshal(bytes, &internetHeader)
```
```gherkin
        Then Unmarshal() should return a nil error
        And I should see struct field values matching the bits in those bytes
```
```go
            log.Println(e == nil)
            // true

            log.Println(internetHeader.RFC791InternetHeaderFormatWord0.Version)
            // 4

            log.Println(internetHeader.RFC791InternetHeaderFormatWord0.IHL)
            // 5
```

See the rest of the [specifications](docs/binary.feature) for error scenarios.

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
