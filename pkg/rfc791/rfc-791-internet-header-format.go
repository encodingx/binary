package rfc791

type RFC791InternetHeaderFormatWithoutOptions struct {
	// Reference: Section 3.1 "Internet Header Format" of
	// RFC 791 Internet Protocol
	// https://datatracker.ietf.org/doc/html/rfc791#section-3.1

	// > A summary of the contents of the internet header follows:
	// >
	// >
	// >   0                   1                   2                   3
	// >   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// >  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// >  |Version|  IHL  |Type of Service|          Total Length         |
	// >  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// >  |         Identification        |Flags|      Fragment Offset    |
	// >  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// >  |  Time to Live |    Protocol   |         Header Checksum       |
	// >  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// >  |                       Source Address                          |
	// >  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// >  |                    Destination Address                        |
	// >  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// >  |                    Options                    |    Padding    |
	// >  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// >
	// >                   Example Internet Datagram Header
	// >
	// >                              Figure 4.
	// >
	// > Note that each tick mark represents one bit position.

	RFC791InternetHeaderFormatWord0 `word:"32"`
	RFC791InternetHeaderFormatWord1 `word:"32"`
	RFC791InternetHeaderFormatWord2 `word:"32"`
	RFC791InternetHeaderFormatWord3 `word:"32"`
	RFC791InternetHeaderFormatWord4 `word:"32"`
}

type RFC791InternetHeaderFormatWord0 struct {
	Version uint8 `bitfield:"4,28"`
	// > Version:  4 bits
	// >
	// >   The Version field indicates the format of the internet header.  This
	// >   document describes version 4.

	IHL uint8 `bitfield:"4,24"`
	// > IHL:  4 bits
	// >
	// >   Internet Header Length is the length of the internet header in 32
	// >   bit words, and thus points to the beginning of the data.  Note that
	// >   the minimum value for a correct header is 5.

	Precedence  uint8 `bitfield:"3,21"`
	Delay       bool  `bitfield:"1,20"`
	Throughput  bool  `bitfield:"1,19"`
	Reliability bool  `bitfield:"1,18"`
	Reserved    uint8 `bitfield:"2,16"`
	// > Type of Service:  8 bits
	// >
	// >   The Type of Service provides an indication of the abstract
	// >   parameters of the quality of service desired.  These parameters are
	// >   to be used to guide the selection of the actual service parameters
	// >   when transmitting a datagram through a particular network.  Several
	// >   networks offer service precedence, which somehow treats high
	// >   precedence traffic as more important than other traffic (generally
	// >   by accepting only traffic above a certain precedence at time of high
	// >   load).  The major choice is a three way tradeoff between low-delay,
	// >   high-reliability, and high-throughput.
	// >
	// >     Bits 0-2:  Precedence.
	// >     Bit    3:  0 = Normal Delay,      1 = Low Delay.
	// >     Bits   4:  0 = Normal Throughput, 1 = High Throughput.
	// >     Bits   5:  0 = Normal Relibility, 1 = High Relibility.
	// >     Bit  6-7:  Reserved for Future Use.
	// >
	// >        0     1     2     3     4     5     6     7
	// >     +-----+-----+-----+-----+-----+-----+-----+-----+
	// >     |                 |     |     |     |     |     |
	// >     |   PRECEDENCE    |  D  |  T  |  R  |  0  |  0  |
	// >     |                 |     |     |     |     |     |
	// >     +-----+-----+-----+-----+-----+-----+-----+-----+
	// >
	// >       Precedence
	// >
	// >         111 - Network Control
	// >         110 - Internetwork Control
	// >         101 - CRITIC/ECP
	// >         100 - Flash Override
	// >         011 - Flash
	// >         010 - Immediate
	// >         001 - Priority
	// >         000 - Routine
	// >
	// >   ...

	TotalLength uint16 `bitfield:"16,0"`
	// > Total Length:  16 bits
	// >
	// >   Total Length is the length of the datagram, measured in octets,
	// >   including internet header and data.  This field allows the length of
	// >   a datagram to be up to 65,535 octets. ...
}

const (
	RFC791InternetHeaderVersion = 4
)

const (
	RFC791InternetHeaderLengthWithoutOptions = 5
)

const (
	RFC791InternetHeaderPrecedenceNetworkControl      = 0b111
	RFC791InternetHeaderPrecedenceInternetworkControl = 0b110
	RFC791InternetHeaderPrecedenceCRITICECP           = 0b101
	RFC791InternetHeaderPrecedenceFlashOverride       = 0b100
	RFC791InternetHeaderPrecedenceFlash               = 0b011
	RFC791InternetHeaderPrecedenceImmediate           = 0b010
	RFC791InternetHeaderPrecedencePriority            = 0b001
	RFC791InternetHeaderPrecedenceRoutine             = 0b000
)

const (
	RFC791InternetHeaderDelayNormal = false
	RFC791InternetHeaderDelayLow    = true
)

const (
	RFC791InternetHeaderThroughputNormal = false
	RFC791InternetHeaderThroughputHigh   = true
)

const (
	RFC791InternetHeaderReliabilityNormal = false
	RFC791InternetHeaderReliabilityHigh   = true
)

type RFC791InternetHeaderFormatWord1 struct {
	Identification uint16 `bitfield:"16,16"`
	// > Identification:  16 bits
	// >
	// >   An identifying value assigned by the sender to aid in assembling the
	// >   fragments of a datagram.

	FlagsBit0Reserved bool `bitfield:"1,15"`
	FlagsBit1         bool `bitfield:"1,14"`
	FlagsBit2         bool `bitfield:"1,13"`
	// > Flags:  3 bits
	// >
	// >   Various Control Flags.
	// >
	// >     Bit 0: reserved, must be zero
	// >     Bit 1: (DF) 0 = May Fragment,  1 = Don't Fragment.
	// >     Bit 2: (MF) 0 = Last Fragment, 1 = More Fragments.
	// >
	// >         0   1   2
	// >       +---+---+---+
	// >       |   | D | M |
	// >       | 0 | F | F |
	// >       +---+---+---+

	FragmentOffset uint16 `bitfield:"13,0"`
	// > Fragment Offset:  13 bits
	// >
	// >   This field indicates where in the datagram this fragment belongs.
	// >   The fragment offset is measured in units of 8 octets (64 bits).  The
	// >   first fragment has offset zero.
}

const (
	RFC791InternetHeaderFlagsBit1MayFragment   = false
	RFC791InternetHeaderFlagsBit1DoNotFragment = true
)

const (
	RFC791InternetHeaderFlagsBit2LastFragment  = false
	RFC791InternetHeaderFlagsBit2MoreFragments = true
)

type RFC791InternetHeaderFormatWord2 struct {
	TimeToLive uint8 `bitfield:"8,24"`
	// > Time to Live:  8 bits
	// >
	// >   This field indicates the maximum time the datagram is allowed to
	// >   remain in the internet system.  If this field contains the value
	// >   zero, then the datagram must be destroyed.  This field is modified
	// >   in internet header processing.  The time is measured in units of
	// >   seconds, but since every module that processes a datagram must
	// >   decrease the TTL by at least one even if it process the datagram in
	// >   less than a second, the TTL must be thought of only as an upper
	// >   bound on the time a datagram may exist.  The intention is to cause
	// >   undeliverable datagrams to be discarded, and to bound the maximum
	// >   datagram lifetime.

	Protocol uint8 `bitfield:"8,16"`
	// > Protocol:  8 bits
	// >
	// >   This field indicates the next level protocol used in the data
	// >   portion of the internet datagram.  The values for various protocols
	// >   are specified in "Assigned Numbers" [9].

	// >                 ASSIGNED INTERNET PROTOCOL NUMBERS
	// >
	// > In the Internet Protocol (IP) [33] there is a field, called Protocol,
	// > to identify the the next level protocol.  This is an 8 bit field.
	// >
	// > Assigned Internet Protocol Numbers
	// >
	// >    Decimal    Octal      Protocol Numbers                  References
	// >    -------    -----      ----------------                  ----------
	// >         0       0         Reserved                              [JBP]
	// >         1       1         ICMP                               [53,JBP]
	// >         2       2         Unassigned                            [JBP]
	// >         3       3         Gateway-to-Gateway              [48,49,VMS]
	// >         4       4         CMCC Gateway Monitoring Message [18,19,DFP]
	// >         5       5         ST                                 [20,JWF]
	// >         6       6         TCP                                [34,JBP]
	// >         7       7         UCL                                    [PK]
	// >         8      10         Unassigned                            [JBP]
	// >         9      11         Secure                                [VGC]
	// >        10      12         BBN RCC Monitoring                    [VMS]
	// >        11      13         NVP                                 [12,DC]
	// >        12      14         PUP                                [4,EAT3]
	// >        13      15         Pluribus                             [RDB2]
	// >        14      16         Telenet                              [RDB2]
	// >        15      17         XNET                              [25,JFH2]
	// >        16      20         Chaos                                [MOON]
	// >        17      21         User Datagram                      [42,JBP]
	// >        18      22         Multiplexing                       [13,JBP]
	// >        19      23         DCN                                  [DLM1]
	// >        20      24         TAC Monitoring                     [55,RH6]
	// >     21-62   25-76         Unassigned                            [JBP]
	// >        63      77         any local network                     [JBP]
	// >        64     100         SATNET and Backroom EXPAK            [DM11]
	// >        65     101         MIT Subnet Support                    [NC3]
	// >     66-68 102-104         Unassigned                            [JBP]
	// >        69     105         SATNET Monitoring                    [DM11]
	// >        70     106         Unassigned                            [JBP]
	// >        71     107         Internet Packet Core Utility         [DM11]
	// >     72-75 110-113         Unassigned                            [JBP]
	// >        76     114         Backroom SATNET Monitoring           [DM11]
	// >        77     115         Unassigned                            [JBP]
	// >        78     116         WIDEBAND Monitoring                  [DM11]
	// >        79     117         WIDEBAND EXPAK                       [DM11]
	// >    80-254 120-376         Unassigned                            [JBP]
	// >       255     377         Reserved                              [JBP]

	HeaderChecksum uint16 `bitfield:"16,0"`
	// > Header Checksum:  16 bits
	// >
	// >   A checksum on the header only.  Since some header fields change
	// >   (e.g., time to live), this is recomputed and verified at each point
	// >   that the internet header is processed.
	// >
	// >   ...
}

const (
	RFC791InternetHeaderProtocolReserved = iota
	RFC791InternetHeaderProtocolICMP
	RFC791InternetHeaderProtocolUnassigned
	RFC791InternetHeaderProtocolGatewayToGateway
	RFC791InternetHeaderProtocolCMCCGatewayMonitoringMessage
	RFC791InternetHeaderProtocolST
	RFC791InternetHeaderProtocolTCP
	RFC791InternetHeaderProtocolUCL
	_
	RFC791InternetHeaderProtocolSecure
	RFC791InternetHeaderProtocolBBNRCCMonitoring
	RFC791InternetHeaderProtocolNVP
	RFC791InternetHeaderProtocolPUP
	RFC791InternetHeaderProtocolPluribus
	RFC791InternetHeaderProtocolTelenet
	RFC791InternetHeaderProtocolXNET
	RFC791InternetHeaderProtocolChaos
	RFC791InternetHeaderProtocolUserDatagram
	RFC791InternetHeaderProtocolMultiplexing
	RFC791InternetHeaderProtocolDCN
	RFC791InternetHeaderProtocolTACMonitoring

	RFC791InternetHeaderProtocolSATNETAndBackroomEXPAK    = 64
	RFC791InternetHeaderProtocolMITSubnetSupport          = 65
	RFC791InternetHeaderProtocolSATNETMonitoring          = 69
	RFC791InternetHeaderProtocolInternetPacketCoreUtility = 71
	RFC791InternetHeaderProtocolBackroomSATNETMonitoring  = 76
	RFC791InternetHeaderProtocolWIDEBANDMonitoring        = 78
	RFC791InternetHeaderProtocolWIDEBANDEXPAK             = 79
)

type RFC791InternetHeaderFormatWord3 struct {
	SourceAddressOctet0 uint8 `bitfield:"8,24"`
	SourceAddressOctet1 uint8 `bitfield:"8,16"`
	SourceAddressOctet2 uint8 `bitfield:"8,8"`
	SourceAddressOctet3 uint8 `bitfield:"8,0"`
	// > Source Address:  32 bits
}

type RFC791InternetHeaderFormatWord4 struct {
	DestinationAddressOctet0 uint8 `bitfield:"8,24"`
	DestinationAddressOctet1 uint8 `bitfield:"8,16"`
	DestinationAddressOctet2 uint8 `bitfield:"8,8"`
	DestinationAddressOctet3 uint8 `bitfield:"8,0"`
	// > Destination Address:  32 bits
}
