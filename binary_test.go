package binary

import (
	"testing"

	"github.com/encodingx/binary/pkg/demo"
)

const (
	destinationAddressOctet = 8
	fragmentOffset          = 8191
	headerChecksum          = 0
	identification          = 0
	sourceAddressOctet      = 1
	timeToLive              = 1
	totalLength             = 65535

	formatBytes  = "%08b"
	formatStruct = "%#v"
)

var (
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

	struct1 demo.RFC791InternetHeaderFormatWithoutOptions

	bytes = []byte{
		0b01000101, 0b11101000, 0b11111111, 0b11111111,
		0b00000000, 0b00000000, 0b01011111, 0b11111111,
		0b00000001, 0b00000110, 0b00000000, 0b00000000,
		0b00000001, 0b00000001, 0b00000001, 0b00000001,
		0b00001000, 0b00001000, 0b00001000, 0b00001000,
	}

	e error
)

func BenchmarkMarshal(b *testing.B) {
	var (
		e error
		i int
	)

	for i = 0; i < b.N; i++ {
		_, e = Marshal(&struct0)
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
		e = Unmarshal(bytes, &struct1)
		if e != nil {
			b.Error(e)
		}
	}
}
