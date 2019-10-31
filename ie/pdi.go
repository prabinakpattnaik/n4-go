package ie

import (
	"fmt"
)

type PDI struct {
	SourceInterface               *InformationElement
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

func NEWPDI(sourceInterface, localFTEID, networkInstance, ueIPAddress, trafficEndpointID, sdffilter, applicaitonID, ethernetPDUSessionInformation, ethernetPacketFilter, qfi, framedRoute, frameRouting, frameIPv6Route *InformationElement) *PDI {
	return &PDI{
		SourceInterface:               sourceInterface,
		LocalFTEID:                    localFTEID,
		NetworkInstance:               networkInstance,
		UEIPAddress:                   ueIPAddress,
		TrafficEndpointID:             trafficEndpointID,
		SDFFilter:                     sdffilter,
		ApplicationID:                 applicaitonID,
		EthernetPDUSessionInformation: ethernetPDUSessionInformation,
		EthernetPacketFilter:          ethernetPacketFilter,
		QFI:                           qfi,
		FramedRoute:                   framedRoute,
		FramedRouting:                 frameRouting,
		FramedIPv6Route:               frameIPv6Route,
	}

}

func (p PDI) Serialize() ([]byte, error) {
	if p.SourceInterface == nil || p.SourceInterface.Type == IEReserved {
		return nil, fmt.Errorf("CreateFAR does not have valid FARID")
	}
	b, err := p.SourceInterface.Serialize()
	if err != nil {
		return nil, fmt.Errorf("FARID serialization error")
	}

	if p.LocalFTEID != nil && p.LocalFTEID.Type != IEReserved {
		b1, err := p.LocalFTEID.Serialize()
		if err != nil {
			return nil, fmt.Errorf("LocalFTEID serialization error")
		}
		b = append(b, b1...)
	}

	if p.NetworkInstance != nil && p.NetworkInstance.Type != IEReserved {
		b1, err := p.NetworkInstance.Serialize()
		if err != nil {
			return nil, fmt.Errorf("NetworkInstance serialization error")
		}
		b = append(b, b1...)
	}

	if p.UEIPAddress != nil && p.UEIPAddress.Type != IEReserved {
		b1, err := p.UEIPAddress.Serialize()
		if err != nil {
			return nil, fmt.Errorf("UEIPAddress serialization error")
		}
		b = append(b, b1...)
	}

	if p.TrafficEndpointID != nil && p.TrafficEndpointID.Type != IEReserved {
		b1, err := p.TrafficEndpointID.Serialize()
		if err != nil {
			return nil, fmt.Errorf("TrafficEndpointID serialization error")
		}
		b = append(b, b1...)
	}

	if p.SDFFilter != nil && p.SDFFilter.Type != IEReserved {
		b1, err := p.SDFFilter.Serialize()
		if err != nil {
			return nil, fmt.Errorf("SDFFilter serialization error")
		}
		b = append(b, b1...)
	}

	if p.ApplicationID != nil && p.ApplicationID.Type != IEReserved {
		b1, err := p.ApplicationID.Serialize()
		if err != nil {
			return nil, fmt.Errorf("ApplicationID serialization error")
		}
		b = append(b, b1...)
	}

	//TODO remaining to be filled

	return b, nil

}
