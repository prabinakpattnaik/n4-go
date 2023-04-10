package ie

import (
	"fmt"

	"github.com/prabinakpattnaik/n4-go/util/util_3gpp"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

type ForwardingParametersIEInFAR struct {
	DestinationInterface *DestinationInterface `tlv:"42"`
	NetworkInstance      *util_3gpp.Dnn        `tlv:"22"`
	RedirectInformation  *RedirectInformation  `tlv:"38"`
	OuterHeaderCreation  *OuterHeaderCreation  `tlv:"84"`
	//TransportLevelMarking   *TransportLevelMarking `tlv:"30"`
	//ForwardingPolicy        *ForwardingPolicy      `tlv:"41"`
	//HeaderEnrichment        *HeaderEnrichment      `tlv:"98"`
	//LinkedTrafficEndpointID *TrafficEndpointID     `tlv:"131"`
	//Proxying                *Proxying              `tlv:"137"`
}

func (f *ForwardingParametersIEInFAR) Serialize() ([]byte, error) {
	if f.DestinationInterface == nil {
		return nil, fmt.Errorf("Mandatory field missing in ForwardfingParameters")
	}
	var b []byte
	d := NewInformationElement(
		IEDestinationInterface,
		0,
		dt.OctetString([]byte{f.DestinationInterface.InterfaceValue}),
	)
	b1, err := d.Serialize()
	if err == nil {
		b = append(b, b1...)
	}

	if f.NetworkInstance != nil {
		b1, err := f.NetworkInstance.Serialize()
		if err == nil {
			networkInstance := NewInformationElement(
				IENetworkInstance,
				0,
				dt.OctetString(b1),
			)
			b1, err = networkInstance.Serialize()
			if err == nil {
				b = append(b, b1...)
			}
		}
	}

	//TODO RedirectInformation. Not included.

	if f.OuterHeaderCreation != nil {
		b1, err = f.OuterHeaderCreation.Serialize()
		if err != nil {
			return nil, err
		}
		ieOHC := NewInformationElement(
			IEOuterHeaderCreation,
			0,
			dt.OctetString(b1),
		)
		b1, err = ieOHC.Serialize()
		if err == nil {
			b = append(b, b1...)
		}
	}
	return b, nil
}
