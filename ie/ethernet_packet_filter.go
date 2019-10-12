package ie

type EthernetPacketFilter struct {
	EthernetFilterID         *InformationElement
	EthernetFilterProperties *InformationElement
	MACAddress               *InformationElement
	Ethertype                *InformationElement
	CTAG                     *InformationElement
	STAG                     *InformationElement
	SDFFilter                *InformationElement
}
