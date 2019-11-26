package ie

import (
	"encoding/binary"
	"net"
)

//OuterHeaderCreation defines struct
type OuterHeaderCreation struct {
	OuterHeaderCreationDescription uint8
	TEID                           uint32
	IPv4Address                    net.IP
	IPv6Address                    net.IP
	PortNumber                     uint16
	//TODO: C-TAG, S-TAG

}

//NewOuterHeaderCreation creates new OuterHeaderCreation
func NewOuterHeaderCreation(ohcd uint8, teid uint32, ipv4address, ipv6address net.IP, pn uint16) *OuterHeaderCreation {
	return &OuterHeaderCreation{
		OuterHeaderCreationDescription: ohcd,
		TEID:                           teid,
		IPv4Address:                    ipv4address,
		IPv6Address:                    ipv6address,
		PortNumber:                     pn,
	}

}

func (ohc OuterHeaderCreation) Serialize() ([]byte, error) {
	b := []byte{ohc.OuterHeaderCreationDescription}
	if ohc.OuterHeaderCreationDescription == 1 || ohc.OuterHeaderCreationDescription == 2 {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, ohc.TEID)
		b = append(b, buf...)
	}

	if ip4 := ohc.IPv4Address.To4(); ip4 != nil {
		b = append(b, ip4...)
	}
	if ip6 := ohc.IPv6Address.To16(); ip6 != nil {
		b = append(b, ip6...)
	}

	if ohc.OuterHeaderCreationDescription == 3 || ohc.OuterHeaderCreationDescription == 4 {
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, ohc.PortNumber)
		b = append(b, buf...)
	}

	return b, nil
}
