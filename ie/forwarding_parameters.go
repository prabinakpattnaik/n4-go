package ie

import (
	"fmt"
)

type ForwardingParameters struct {
	DestinationInterface    *InformationElement
	NetworkInstance         *InformationElement
	RedirectInformation     *InformationElement
	OuterHeaderCreation     *InformationElement
	TransportLevelMarking   *InformationElement
	ForwardingPolicy        *InformationElement
	HeaderEnrichment        *InformationElement
	LinkedTrafficEndpointID *InformationElement
	Proxying                *InformationElement
}

func NewForwardingParameters(d, n, r, o, t, f, h, l, p *InformationElement) *ForwardingParameters {
	return &ForwardingParameters{
		DestinationInterface:    d,
		NetworkInstance:         n,
		RedirectInformation:     r,
		OuterHeaderCreation:     o,
		TransportLevelMarking:   t,
		ForwardingPolicy:        f,
		HeaderEnrichment:        h,
		LinkedTrafficEndpointID: l,
		Proxying:                p,
	}

}

func (f ForwardingParameters) Serialize() ([]byte, error) {
	if f.DestinationInterface == nil || f.DestinationInterface.Type == IEReserved {
		return nil, fmt.Errorf("ForwardingParameters does not have valid DestionationInterface")
	}
	b, err := f.DestinationInterface.Serialize()
	if err != nil {
		return nil, fmt.Errorf("DestinationInterface serialization error")
	}

	if f.NetworkInstance != nil && f.NetworkInstance.Type != IEReserved {
		b1, err := f.NetworkInstance.Serialize()
		if err != nil {
			return nil, fmt.Errorf("NetworkInstance serialization error")
		}
		b = append(b, b1...)
	}

	if f.RedirectInformation != nil && f.RedirectInformation.Type != IEReserved {
		b1, err := f.RedirectInformation.Serialize()
		if err != nil {
			return nil, fmt.Errorf("RedirectInformation serialization error")
		}
		b = append(b, b1...)
	}

	if f.OuterHeaderCreation != nil && f.OuterHeaderCreation.Type != IEReserved {
		b1, err := f.OuterHeaderCreation.Serialize()
		if err != nil {
			return nil, fmt.Errorf("OuterHeaderCreation serialization error")
		}
		b = append(b, b1...)
	}

	if f.TransportLevelMarking != nil && f.TransportLevelMarking.Type != IEReserved {
		b1, err := f.TransportLevelMarking.Serialize()
		if err != nil {
			return nil, fmt.Errorf("TransportLevelMarking serialization error")
		}
		b = append(b, b1...)
	}

	if f.ForwardingPolicy != nil && f.ForwardingPolicy.Type != IEReserved {
		b1, err := f.ForwardingPolicy.Serialize()
		if err != nil {
			return nil, fmt.Errorf("ForwardingPolicy serialization error")
		}
		b = append(b, b1...)
	}

	if f.HeaderEnrichment != nil && f.HeaderEnrichment.Type != IEReserved {
		b1, err := f.HeaderEnrichment.Serialize()
		if err != nil {
			return nil, fmt.Errorf("HeaderEnrichment serialization error")
		}
		b = append(b, b1...)
	}

	if f.LinkedTrafficEndpointID != nil && f.LinkedTrafficEndpointID.Type != IEReserved {
		b1, err := f.LinkedTrafficEndpointID.Serialize()
		if err != nil {
			return nil, fmt.Errorf("LinkedTrafficEndpointID serialization error")
		}
		b = append(b, b1...)
	}

	if f.Proxying != nil && f.Proxying.Type != IEReserved {
		b1, err := f.Proxying.Serialize()
		if err != nil {
			return nil, fmt.Errorf("Proxying serialization error")
		}
		b = append(b, b1...)
	}

	return b, nil

}
