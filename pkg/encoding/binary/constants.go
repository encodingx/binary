package binary

const (
	wordLengthInBytes = 4
	wordLengthInBits  = wordLengthInBytes * 8
	wordRangeMaximum  = (1 << wordLengthInBits) - 1
)
