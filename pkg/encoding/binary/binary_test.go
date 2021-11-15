package binary

import (
	"testing"
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

func BenchmarkMarshal(b *testing.B) {
	var (
		e error
		i int
	)

	for i = 0; i < b.N; i++ {
		_, e = Marshal(structure0)
		if e != nil {
			b.Error(e)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	var (
		e error
		i int
	)

	for i = 0; i < b.N; i++ {
		e = Unmarshal(binary, &structure1)
		if e != nil {
			return
		}
	}
}
