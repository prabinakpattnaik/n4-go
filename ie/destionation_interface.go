package ie

const (
	DestinationInterfaceAccess uint8 = iota
	DestinationInterfaceCore
	DestinationInterfaceSgiLanN6Lan
	DestinationInterfaceCpFunction
	DestinationInterfaceLiFunction
)

type DestinationInterface struct {
	InterfaceValue uint8 // 0x00001111
}

func (d *DestinationInterface) Serialize() ([]byte, error) {
	return []byte{d.InterfaceValue}, nil
}
