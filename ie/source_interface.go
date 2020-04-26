package ie

const (
	SourceInterfaceAccess uint8 = iota
	SourceInterfaceCore
	SourceInterfaceSgiLanN6Lan
	SourceInterfaceCpFunction
)

type SourceInterface struct {
	InterfaceValue uint8 // 0x00001111
}

func (s *SourceInterface) Serialize() ([]byte, error) {
	var b []byte
	b[0] = s.InterfaceValue
	return b, nil
}
