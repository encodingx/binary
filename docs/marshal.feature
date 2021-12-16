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
        And the format is a series of "bit fields"
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
        And a format is represented by a type definition of a "format-struct"
        And the format-struct nests one or more exported "word-structs"
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
        And each word-struct has exported field(s) corresponding to bit field(s)
        And the fields are of unsigned integer or boolean types
        And the fields are tagged to indicate the lengths of those bit fields
        And the length of each bit field does not overflow the type of the field
            """
            A bit field overflows a type
            when it is long enough to represent values
            outside the set of values of the type.
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

    Scenario:
        Given a variable that is not a pointer
        When I pass to Marshal() as an argument such a variable
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            Argument to Marshal() should be a pointer to a format-struct.
            Argument to Marshal() is not a pointer.
            """

    Scenario:
        Given a pointer that does not point to a struct variable
        When I pass to Marshal() such a pointer
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            Argument to Marshal() should be a pointer to a format-struct.
            Argument to Marshal() does not point to a struct variable.
            """

    Scenario:
        Given a format-struct with no exported fields
        When I pass to Marshal() a pointer to such a format-struct
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            A format-struct should nest exported word-structs.
            Argument to Marshal() points to a format-struct '<FormatStructType>'
            with no exported fields.
            """

    Scenario:
        Given an exported field in a format-struct is not of type struct
        When I pass to Marshal() a pointer to such a format-struct
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            A format-struct should nest exported word-structs.
            Argument to Marshal() points to a format-struct '<FormatStructType>'
            with an exported field '<NameOfStructField>' that is not a struct.
            """

    Scenario:
        Given an exported field in a format-struct with no struct tag
        When I pass to Marshal() a pointer to such a format-struct
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            Exported fields in a format-struct should be tagged
            with a key "word" and a value
            indicating the length of the word in number of bits
            (e.g. `word:"32"`).
            Argument to Marshal() points to a format-struct '<FormatStructType>'
            with an exported field '<NameOfStructField>' that has no struct tag.
            """

    Scenario:
        Given an exported field in a format-struct with a malformed struct tag
            """
            A struct tag is malformed when its key is not "word"
            or when its value cannot be parsed as an unsigned integer.
            """
        When I pass to Marshal() a pointer to such a format-struct
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            Exported fields in a format-struct should be tagged
            with a key "word" and a value
            indicating the length of the word in number of bits
            (e.g. `word:"32"`).
            Argument to Marshal() points to a format-struct '<FormatStructType>'
            with an exported field '<NameOfStructField>'
            that has a malformed struct tag: <message of wrapped error>.
            """

    Scenario:
        Given a word of length not a multiple of eight in the range [8, 64]
        When I pass to Marshal() a pointer to a format-struct nesting such word
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            The length of a word should be a multiple of eight
            in the range [8, 64].
            Argument to Marshal() points to a format-struct '<FormatStructType>'
            containing a word-struct '<NameOfStructField>' of length <length>
            not in {8, 16, 24, ... 64}.
            """

    Scenario:
        Given a word-struct containing no exported fields
        When I pass to Marshal() a pointer to a format-struct nesting such word
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            A word-struct should nest exported fields representing bit fields.
            Argument to Marshal() points to a format-struct '<FormatStructType>'
            containing a word-struct '<FieldName>' of type '<WordStructType>',
            which has no exported fields.
            """

    Scenario:
        Given a word-struct containing a field of unsupported type
            """
            Supported types are uint, uintN where N = {8, 16, 32, 64} and bool.
            """
        When I pass to Marshal() a pointer to a format-struct nesting such word
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            The fields of a word-struct should be of type uintN or bool.
            Argument to Marshal() points to a format-struct '<FormatStructType>'
            containing a word-struct '<FieldName0>' of type '<WordStructType>',
            which has a field '<FieldName1>' of unsupported type '<FieldType>'.
            """

    Scenario:
        Given a word-struct containing a field with no struct tag
        When I pass to Marshal() a pointer to a format-struct nesting such word
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            Exported fields in a word-struct should be tagged
            with a key "bitfield" and a value
            indicating the length of the bit field in number of bits
            (e.g. `bitfield:"1"`).
            Argument to Marshal() points to a format-struct '<FormatStructType>'
            containing a word-struct '<FieldName0>' of type '<WordStructType>',
            which has a field '<FieldName1>' that has no struct tag.
            """

    Scenario:
        Given a word-struct containing a field with a malformed struct tag
            """
            A struct tag is malformed when its key is not "bitfield"
            or when its value cannot be parsed as an unsigned integer.
            """
        When I pass to Marshal() a pointer to a format-struct nesting such word
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            Exported fields in a word-struct should be tagged
            with a key "bitfield" and a value
            indicating the length of the bit field in number of bits
            (e.g. `bitfield:"1"`).
            Argument to Marshal() points to a format-struct '<FormatStructType>'
            containing a word-struct '<FieldName0>' of type '<WordStructType>',
            which has a field '<FieldName1>' that has a malformed struct tag:
            <message of wrapped error>.
            """

    Scenario:
        Given a word-struct with a bit field of length overflowing its type
            """
            A bit field overflows a type
            when it is long enough to represent values
            outside the set of values of the type.
            """
        When I pass to Marshal() a pointer to a format-struct nesting such word
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            The number of unique values a bit field can contain
            must not exceed the size of its type.
            Argument to Marshal() points to a format-struct '<FormatStructType>'
            containing a word-struct '<FieldName0>' of type '<WordStructType>',
            which has a bit field '<FieldName1>' of length <length>
            exceeding the size of type <FieldType>, <size>.
            """

    Scenario:
        Given a word of length not equal to the sum of lengths of its bit fields
        When I pass to Marshal() a pointer to a format-struct nesting such word
        Then Marshal() should return a byte slice of zero length and an error
            """
            Marshal error:
            The length of a word
            should be equal to the sum of lengths of its bit fields.
            Argument to Marshal() points to a format-struct '<NameOfStructType>'
            containing a word-struct '<NameOfStructField>' of length <length>
            not equal to the sum of the lengths of its bit fields, <sum>.
            """