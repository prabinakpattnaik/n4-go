package msg

import (
	"github.com/prabinakpattnaik/n4-go/ie"
)

//PFCPSessionEstablishmentResponse struct defines the structure
type PFCPSessionEstablishmentResponse struct {
	Header                     *PFCPHeader
	NodeID                     *ie.InformationElement
	Cause                      *ie.InformationElement
	OffendingIE                *ie.InformationElement
	UPFSEID                    *ie.InformationElement
	CreatedPDR                 *ie.InformationElement
	LoadControlInformation     *ie.InformationElement
	OverloadControlInformation *ie.InformationElement
	FailedRuleID               *ie.InformationElement
	CreatedTrafficEndpoint     *ie.InformationElement
}

//NewPFCPSessionEstablishmentResponse creates new PFCPSessionEstablishmentResponse object
func NewPFCPSessionEstablishmentResponse(h *PFCPHeader, n, c, offending, u, createdPDR, l, overload, f, cte *ie.InformationElement) PFCPSessionEstablishmentResponse {
	return PFCPSessionEstablishmentResponse{
		Header:                     h,
		NodeID:                     n,
		Cause:                      c,
		OffendingIE:                offending,
		UPFSEID:                    u,
		CreatedPDR:                 createdPDR,
		LoadControlInformation:     l,
		OverloadControlInformation: overload,
		FailedRuleID:               f,
		CreatedTrafficEndpoint:     cte,
	}

}

//Serialize function converts PFCPSessionEstablishmentResponse into byte array
func (sr PFCPSessionEstablishmentResponse) Serialize() ([]byte, error) {
	var b []byte
	if sr.NodeID == nil || sr.Cause == nil {
		return b, nil
	}

	dataLength := sr.Len()
	output := make([]byte, dataLength)
	pfcpend := uint16(PFCPBasicHeaderLength) + PFCPMessageSize
	copy(output[:pfcpend], sr.Header.Serialize())

	nb, err := sr.NodeID.Serialize()
	if err != nil {
		return nil, err
	}
	nodeIDEnd := pfcpend + ie.IEBasicHeaderSize + sr.NodeID.Len()
	copy(output[pfcpend:nodeIDEnd], nb)

	cause, err := sr.Cause.Serialize()
	if err != nil {
		return nil, err
	}
	causeEnd := nodeIDEnd + ie.IEBasicHeaderSize + sr.Cause.Len()
	copy(output[nodeIDEnd:causeEnd], cause)

	var offendingIEEnd, upfseidEnd, createdPDREnd uint16

	if sr.OffendingIE != nil {
		offendingIEEnd = causeEnd + ie.IEBasicHeaderSize + sr.OffendingIE.Len()
		o, err := sr.OffendingIE.Serialize()
		if err != nil {
			return nil, err
		}
		copy(output[causeEnd:offendingIEEnd], o)

	}

	//TODO: it applies only when cause is success
	if sr.UPFSEID != nil && offendingIEEnd == 0 {
		upfseidEnd = causeEnd + ie.IEBasicHeaderSize + sr.UPFSEID.Len()
		u, err := sr.UPFSEID.Serialize()
		if err != nil {
			return nil, err
		}
		copy(output[causeEnd:upfseidEnd], u)

	}

	if offendingIEEnd == 0 && upfseidEnd > 0 && sr.CreatedPDR != nil {
		createdPDREnd = upfseidEnd + ie.IEBasicHeaderSize + sr.CreatedPDR.Len()
		c, err := sr.CreatedPDR.Serialize()
		if err != nil {
			return nil, err
		}
		copy(output[upfseidEnd:createdPDREnd], c)
	}
	//TODO: remaining to be added (load control information,OverloadControlInformation,FailedRuleID, and CreatedTrafficEndpoint )

	return output, nil

}

func (sr PFCPSessionEstablishmentResponse) Len() uint16 {
	return uint16(PFCPBasicHeaderLength) + sr.Header.MessageLength
}

func (sr PFCPSessionEstablishmentResponse) Type() PFCPType {
	return sr.Header.MessageType
}

func (sr PFCPSessionEstablishmentResponse) GetHeader() *PFCPHeader {
	return sr.Header
}
