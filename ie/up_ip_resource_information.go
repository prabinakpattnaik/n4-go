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
	SourceInterface uint8
}

func NewUPIPResourceInformation(v4, v6 bool, teidri uint8, assoni, assosi bool, teidrange uint8, ipv4address, ipv6address net.IP, networkinstance []byte, sourceinterface uint8) *UPIPResourceInformation {
	return &UPIPResourceInformation{
		V4:              v4,
		V6:              v6,
		TEIDRI:          teidri,
		ASSONI:          assoni,
		ASSOSI:          assosi,
		TEIDRange:       teidrange,
		IPv4Address:     ipv4address,
		IPv6Address:     ipv6address,
		NetworkInstance: networkinstance,
		SourceInterface: sourceinterface,
	}

}

func NewUPIPResourceInformationFromByte(length uint8, input []byte) *UPIPResourceInformation {
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
		V4:          isV4,
		V6:          isV6,
		TEIDRI:      teidRi,
		ASSONI:      assoni,
		ASSOSI:      assosi,
		IPv4Address: ip4address,
		IPv6Address: ip6address,
	}

}
func (u UPIPResourceInformation) Serialize() ([]byte, error) {
	var firstByte uint8
	if u.V4 {
		firstByte = 1
	}
	if u.V6 {
		firstByte |= 2
	}
	//TODO three bit is maximum size
	firstByte |= (u.TEIDRI << 2)
	if u.ASSONI {
		firstByte |= 32
	}
	if u.ASSOSI {
		firstByte |= 64
	}

	b := make([]byte, 2)
	b[0] = firstByte
	b[1] = u.TEIDRange

	if u.V4 {
		b = append(b, u.IPv4Address...)
	}
	if u.V6 {
		b = append(b, u.IPv6Address...)
	}

	if u.ASSONI {
		b = append(b, u.NetworkInstance...)
	}
	if u.ASSOSI {
		b = append(b, u.SourceInterface)
	}

	return b, nil

}
