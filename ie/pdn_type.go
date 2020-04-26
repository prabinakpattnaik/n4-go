package ie

const (
	PDNTypeIpv4 uint8 = iota + 1
	PDNTypeIpv6
	PDNTypeIpv4v6
	PDNTypeNonIp
	PDNTypeEthernet
)

type PDNType struct {
	PdnType uint8 // 0x00000111
}

func (p *PDNType) Serialize() (data []byte, err error) {
    return []byte{p.PdnType}, nil
}

