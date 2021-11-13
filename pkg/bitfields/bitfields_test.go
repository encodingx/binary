package bitfields

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type rfc791InternetHeaderWord0 struct {
	Version                  uint `bitfield:"4,28"`
	InternetHeaderLength     uint `bitfield:"4,24"`
	TypeOfServicePrecedence  uint `bitfield:"3,21"`
	TypeOfServiceDelay       bool `bitfield:"1,20"`
	TypeOfServiceThroughput  bool `bitfield:"1,19"`
	TypeOfServiceReliability bool `bitfield:"1,18"`
	TypeOfServiceReserved    uint `bitfield:"2,16"`
	TotalLength              uint `bitfield:"16,0"`
}

func TestBitfields(t *testing.T) {
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
		internetHeaderPrecedenceRoutine = iota
		internetHeaderPrecedencePriority
		internetHeaderPrecedenceImmediate
		internetHeaderPrecedenceFlash
		internetHeaderPrecedenceFlashOverride
		internetHeaderPrecedenceCRITICOrECP
		internetHeaderPrecedenceInternetworkControl
		internetHeaderPrecedenceNetworkControl
	)

	var (
		binary = []byte{0b01000101, 0b01010100, 0b1110110, 0b00001110}

		structure0 = rfc791InternetHeaderWord0{
			Version:                  internetHeaderVersion,
			InternetHeaderLength:     internetHeaderLength,
			TypeOfServicePrecedence:  internetHeaderPrecedenceImmediate,
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

	bytes, e = Marshal32BitWord(structure0)
	if e != nil {
		t.Error(e)
	}

	assert.Equal(t,
		binary, bytes,
		"The slice of bytes returned by function Marshal32BitWord "+
			"should match the expected given the struct field values.",
	)

	e = Unmarshal32BitWord(&structure1, binary)
	if e != nil {
		return
	}

	assert.Equal(t,
		structure0, structure1,
		"A struct marshalled by function Marshal32BitWord and "+
			"then unmarshalled by function Unmarshal32BitWord "+
			"should by the same as the original struct.",
	)
}
