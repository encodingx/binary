# Version 2

Feature: Marshal

    As a Go developer implementing a binary message or file format,
    I want a function that converts a struct into a series of bits
    so that I can avoid the complexities of custom bit manipulation.

    Background:
        Given a message or file "format"
            """
            A format specifies how bits are arranged to encode information.
            """
        And this format is a series of "bit fields"
            """
            A bit field is one or more adjacent bits representing a value,
            and should not be confused with struct fields.
            """
        And adjacent bit fields are grouped into "words"
            """
            A word is a series of bits that can be simultaneously processed.
            The length of a word is limited by computer architecture
            and programming language design (64 bits in Go).
            """
        And a format is represented by a struct type
        And the struct type nests exported struct type(s) representing word(s)
        And the words are tagged to indicate their lengths in number of bits
        And the length of each word is a multiple of eight in the range [8, 64]
            """
            type RFC791InternetHeaderFormatWithoutOptions struct {
                RFC791InternetHeaderFormatWord0 `word:"32"`
                RFC791InternetHeaderFormatWord1 `word:"32"`
                RFC791InternetHeaderFormatWord2 `word:"32"`
                RFC791InternetHeaderFormatWord3 `word:"32"`
                RFC791InternetHeaderFormatWord4 `word:"32"`
            }
            """
        And each word struct has exported field(s) representing bit field(s)
        And the fields are of unsigned integer or boolean types
        And the fields are tagged to indicate their lengths in number of bits
        And the length of each field does not overflow the type of the field
            """
            A field overflows its type when it has bits to represent values
            outside the set of values that can be represented by that type.
            """
        And the sum of lengths of all fields is equal to the length of that word
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

    Scenario:
        Given a proper struct variable representing a binary message or file
            """
            // See Background for a definition of a "proper" struct.
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
            when it exceeds the range of values
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
        Then Marshal() should return the message/file as a slice of bytes
        And I should see struct field values reflected as bits in those bytes
            """
            log.Printf("%08b", bytes)
            // [01000101 ...]
            """
        And Marshal() should return a nil error alongside that byte slice
            """
            log.Println(e == nil)
            // true
            """
