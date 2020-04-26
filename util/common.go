package util

const (
	DestinationInterfaceAccess uint8 = iota
	DestinationInterfaceCore
	DestinationInterfaceSgiLanN6Lan
	DestinationInterfaceCpFunction
	DestinationInterfaceLiFunction
)
const (
	SourceInterfaceAccess uint8 = iota
	SourceInterfaceCore
	SourceInterfaceSgiLanN6Lan
	SourceInterfaceCpFunction
)
