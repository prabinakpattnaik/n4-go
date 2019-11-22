package ie

import "net"

type UEIPAddress struct {
	V6                       bool
	V4                       bool
	SD                       bool
	IPV6D                    bool
	IPv4Address              net.IP
	IPv6Address              net.IP
	IPv6PrefixDelegationBits uint8
}

func NewUEIPAddress(v6, v4, sd, ipv6d bool, ipv4address, ipv6address net.IP, ipv6PDB uint8) *UEIPAddress {
	return &UEIPAddress{
		V6:                       v6,
		V4:                       v4,
		SD:                       sd,
		IPV6D:                    ipv6d,
		IPv4Address:              ipv4address,
		IPv6Address:              ipv6address,
		IPv6PrefixDelegationBits: ipv6PDB,
	}

}

func (u *UEIPAddress) Serialize() ([]byte, error) {
	var b, b1 byte
	var ip []byte
	//UEIPaddress can be either IPv4 or IPv6
	if u.V6 {
		b = 1
		ip = u.IPv6Address
	}
	if u.V4 {
		b = 2
		ip = u.IPv4Address.To4()
	}
	if u.SD {
		b = b | 4
	}
	if u.IPV6D {
		b = b | 8
		b1 = u.IPv6PrefixDelegationBits
	}
	f := []byte{b}
	f = append(f, ip...)
	if u.IPV6D {
		f = append(f, b1)
	}
	return f, nil
}
