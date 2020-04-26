package ie

import (
	"fmt"

	"bitbucket.org/sothy5/n4-go/util/util_3gpp"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

type PDI struct {
	SourceInterface *SourceInterface `tlv:"20"`
	LocalFTEID      *FTEID           `tlv:"21"`
	NetworkInstance *util_3gpp.Dnn   `tlv:"22"`
	UEIPAddress     *UEIPAddress     `tlv:"93"`
	//TrafficEndpointID             *TrafficEndpointID             `tlv:"131"`
	SDFFilter       *SDFFilter                     `tlv:"23"`
	//ApplicationID                 *ApplicationID                 `tlv:"24"`
	//EthernetPDUSessionInformation *EthernetPDUSessionInformation `tlv:"142"`
	//EthernetPacketFilter          *EthernetPacketFilter          `tlv:"132"`
	//QFI                           *QFI                           `tlv:"124"`
	//FramedRoute                   *FramedRoute                   `tlv:"153"`
	//FramedRouting                 *FramedRouting                 `tlv:"154"`
	//FramedIPv6Route               *FramedIPv6Route               `tlv:"155"`
}

func (p *PDI) Serialize() ([]byte, error) {
	if p.SourceInterface == nil {
		return nil, fmt.Errorf("No valid PDRID")
	}

	si := NewInformationElement(
		IESourceInterface,
		0,
		dt.OctetString(p.SourceInterface.InterfaceValue),
	)
	b, _ := si.Serialize()

	if p.LocalFTEID != nil {
		b1, err := p.LocalFTEID.Serialize()
		if err == nil {
			lFTEID := NewInformationElement(
				IEFTEID,
				0,
				dt.OctetString(b1),
			)
			b1, err := lFTEID.Serialize()
			if err == nil {
				b = append(b, b1...)
			}
		}
	}

	if p.NetworkInstance != nil {
		b1, err := p.NetworkInstance.Serialize()
		if err == nil {
			networkInstance := NewInformationElement(
				IENetworkInstance,
				0,
				dt.OctetString(b1),
			)
			b1, err = networkInstance.Serialize()
			fmt.Printf("data [%x]\n", b1)
			if err == nil {
				b = append(b, b1...)
			}
		}
	}

	if p.UEIPAddress != nil {
		b1, err := p.UEIPAddress.Serialize()
		if err == nil {
			ueIPaddress := NewInformationElement(
				IEUEIPaddress,
				0,
				dt.OctetString(b1),
			)
			b1, err = ueIPaddress.Serialize()
			if err == nil {
				b = append(b, b1...)
			}
		}
	}

	if p.SDFFilter !=nil {
	    b1, err := p.SDFFilter.Serialize()
	    if err != nil{
	    return nil, err
	    }
	    sdfFilter := NewInformationElement(
	            IESDFFilter,
		    0,
		    dt.OctetString(b1),
	    )
	    b1,err = sdfFilter.Serialize()
	    if err != nil {
	    return nil, err
	    }
	    b=append(b,b1...)
	}

	return b, nil
}
