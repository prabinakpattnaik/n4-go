package ie

import (
	"net"
)

type UPIPResourceInformation struct {
	V4              bool
	V6              bool
	TEIDRI          uint8
	ASSONI          bool
	ASSOSI          bool
	TEIDRange       uint8
	IPv4Address     net.IP
	IPv6Address     net.IP
	NetworkInstance []byte
	SourceInterface []byte
}

func NewUPIPResourceInformation(length uint8, input []byte) *UPIPResourceInformation {
	if length == 0 {
		return nil
	}

	firstByte := input[0]
	isV4 := (uint8(firstByte&0x01) == 1)
	isV6 := (uint8(firstByte&0x02) == 1)
	teidRi := firstByte & 0x1C
	assoni := (uint8(firstByte&0x20) == 1)
	assosi := (uint8(firstByte&0x40) == 1)

	//eighth bit of firstbyte is not tested
	if !(isV4 || isV6) {
		return nil
	}
	//importance of assoni and assosi is not taken
	// TEID range indication

	var ip4address, ip6address []byte

	if isV4 {
		ip4address = input[2:6]
		if isV6 {
			ip6address = input[6:21]
		}
	} else {
		if isV6 {
			ip6address = input[2:17]
		}
	}
	//Network instance and source interface details missing

	return &UPIPResourceInformation{
		V4:          bool(isV4),
		V6:          bool(isV6),
		TEIDRI:      teidRi,
		ASSONI:      assoni,
		ASSOSI:      assosi,
		IPv4Address: ip4address,
		IPv6Address: ip6address,
	}

}
