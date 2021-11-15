package codecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestCodec(t *testing.T) {
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

	var (
		binary = []byte{0b01000101, 0b01010100, 0b1110110, 0b00001110}

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
		codec Codec
		e     error
	)

	codec = NewCodec()

	bytes, e = codec.Marshal(&structure0)
	if e != nil {
		t.Error(e)
	}

	assert.Equal(t,
		binary, bytes,
		"The slice of bytes returned by function Codec.Marshal "+
			"should match the expected given the struct field values.",
	)

	e = codec.Unmarshal(binary, &structure1)
	if e != nil {
		t.Error(e)
	}

	assert.Equal(t,
		structure0, structure1,
		"A struct marshalled by function Codec.Marshal and "+
			"then unmarshalled by function Codec.Unmarshal "+
			"should by the same as the original struct.",
	)
}
