package ie

import (
	"encoding/binary"
	"net"
)

//UP Function Features is a struct
//Table 8.2.25-1
type FSEID struct {
	V4         bool
	V6         bool
	SEID       uint64
	IP4Address net.IP
	IP6Address net.IP
}

func NewFSEID(v4, v6 bool, seid uint64, ip4address, ip6address net.IP) *FSEID {

	return &FSEID{
		V4:         v4,
		V6:         v6,
		SEID:       seid,
		IP4Address: ip4address,
		IP6Address: ip6address,
	}
}

func (f *FSEID) Serialize() ([]byte, error) {
	var firstbyte uint8
	var b []byte
	var ip []byte
	if f.V6 {
		firstbyte = uint8(1)
		ip = f.IP4Address.To16()
	}
	if f.V4 {
		firstbyte = uint8(2)
		ip = f.IP4Address.To4()
	}
	seid := make([]byte, 8)
	binary.BigEndian.PutUint64(seid, uint64(f.SEID))
	b = append(b, firstbyte)
	b = append(b, seid...)
	b = append(b, ip...)
	return b, nil
}

func NewFSEIDFromByte(data []byte) *FSEID {
	v6 := (data[0]&0x01 == 1)
	v4 := (data[0]&0x02 == 2)
	seid := binary.BigEndian.Uint64(data[1:9])
	var ip4address, ip6address net.IP
	if v4 && v6 {

		ip4address = data[9:12]
		ip6address = data[13:]
	} else if v4 {
		ip4address = data[9:]
	} else if v6 {
		ip6address = data[9:]

	}

	return NewFSEID(v4, v6, seid, ip4address, ip6address)
}
