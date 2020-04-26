package ie

const (
	OuterHeaderRemovalGtpUUdpIpv4 uint8 = iota
	OuterHeaderRemovalGtpUUdpIpv6
	OuterHeaderRemovalUdpIpv4
	OuterHeaderRemovalUdpIpv6
)

type OuterHeaderRemoval struct {
	OuterHeaderRemovalDescription uint8
}

func(o *OuterHeaderRemoval)Serialize() ([]byte, error){
     var b []byte
     b[0]=o.OuterHeaderRemovalDescription
     return b, nil
}
