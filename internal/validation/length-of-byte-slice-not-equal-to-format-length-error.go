package validation

import (
	"fmt"
)

func NewLengthOfByteSliceNotEqualToFormatLengthError(
	formatName string, formatLength, byteSliceLength uint,
) (
	e error,
) {
	const (
		format = "" +
			"A byte slice into which a format-struct would be unmarshalled " +
			"should be of length equal to the sum of lengths of words " +
			"in the format represented by the struct. " +
			"Argument to Unmarshal points to a format-struct \"%s\" " +
			"of length %d bits " +
			"not equal to the length of the byte slice, %d bits."
	)

	return fmt.Errorf(format, formatName, formatLength, byteSliceLength)
}
