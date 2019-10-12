package ie

type PDI struct {
	LocalFTEID                    *InformationElement
	NetworkInstance               *InformationElement
	UEIPAddress                   *InformationElement
	TrafficEndpointID             *InformationElement
	SDFFilter                     *InformationElement
	ApplicationID                 *InformationElement
	EthernetPDUSessionInformation *InformationElement
	EthernetPacketFilter          *InformationElement
	QFI                           *InformationElement
	FramedRoute                   *InformationElement
	FramedRouting                 *InformationElement
	FramedIPv6Route               *InformationElement
}
