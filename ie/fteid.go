package ie

import (
	"encoding/binary"
	"net"
)

type FTEID struct {
	IsV4        bool
	IsV6        bool
	CH          bool
	CHID        bool
	TEID        uint32
	IPv4Address net.IP
	IPv6Address net.IP
	CHOOSEID    uint8
}

//NewFTEID creates new FTEID struct
func NewFTEID(isv4, isv6, ch, chid bool, teid uint32, ipv4address, ipv6address net.IP, chooseid uint8) *FTEID {
	return &FTEID{
		IsV4:        isv4,
		IsV6:        isv6,
		CH:          ch,
		CHID:        chid,
		TEID:        teid,
		IPv4Address: ipv4address,
		IPv6Address: ipv6address,
		CHOOSEID:    chooseid,
	}

}

//Serialize converts FTEID into byte array
func (f *FTEID) Serialize() ([]byte, error) {
	var firstbyte uint8
	var b []byte
	var ip []byte

	if f.IsV4 {
		firstbyte = uint8(1)

	}
	if f.IsV6 {
		firstbyte |= uint8(2)

	}
	if f.CH {
		firstbyte |= uint8(4)

	}
	if f.CHID {
		firstbyte |= uint8(16)
	}
	fteid := []byte{firstbyte}
	if !f.CH {
		teid := make([]byte, 4)
		binary.BigEndian.PutUint32(teid, uint32(f.TEID))
		fteid = append(fteid, teid...)
	}

	if !f.CH && f.IsV4 {
		ip = f.IPv4Address.To4()
		fteid = append(fteid, ip...)

	}
	if !f.CH && f.IsV6 {

		fteid = append(fteid, f.IPv6Address.To16()...)

	}

	if f.CHID {
		fteid = append(fteid, f.CHOOSEID)
	}

	return b, nil
}
