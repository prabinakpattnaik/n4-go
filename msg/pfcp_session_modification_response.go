package msg

import (
	"fmt"

	"bitbucket.org/sothy5/n4-go/ie"
)

type PFCPSessionModificationResponse struct {
	Header                     *PFCPHeader
	Cause                      *ie.InformationElement
	OffendingIE                *ie.InformationElement
	CreatedPDR                 *ie.InformationElement
	LoadControlInformation     *ie.InformationElement
	OverloadControlInformation *ie.InformationElement
	UsageReport                *ie.InformationElement
	FailedRuleID               *ie.InformationElement
	AdditionalUsageRI          *ie.InformationElement
	CreatedTrafficEndpoint     *ie.InformationElement
}

func NewPFCPSessionModificationResponse(h *PFCPHeader, c, offending, createdPDR *ie.InformationElement) PFCPSessionModificationResponse {
	return PFCPSessionModificationResponse{
		Header:      h,
		Cause:       c,
		OffendingIE: offending,
		CreatedPDR:  createdPDR,
	}
}

//Serialize function converts PFCPSessionModificationResponse into byte array
func (smr PFCPSessionModificationResponse) Serialize() ([]byte, error) {
	var b []byte
	if smr.Cause == nil {
		return nil, fmt.Errorf("There is no Cause value")
	}

	dataLength := smr.Len()
	fmt.Printf("data Legnth %d\n", dataLength)
	output := make([]byte, dataLength)
	pfcpend := uint16(PFCPBasicHeaderLength) + PFCPMessageSize
	copy(output[:pfcpend], smr.Header.Serialize())

	b, err := smr.Cause.Serialize()
	if err != nil {
		return nil, err
	}

	if smr.OffendingIE != nil {
		o, err := smr.OffendingIE.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, o...)

	}

	if smr.CreatedPDR != nil {
		c, err := smr.CreatedPDR.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, c...)
	}

	copy(output[pfcpend:], b)
	//TODO: remaining to be added (load control information,OverloadControlInformation, UsageReport,FailedRuleID, and CreatedTrafficEndpoint )

	return output, nil

}

func (sr PFCPSessionModificationResponse) Len() uint16 {
	return uint16(PFCPBasicHeaderLength) + sr.Header.MessageLength
}

func (sr PFCPSessionModificationResponse) Type() PFCPType {
	return sr.Header.MessageType
}

func (sr PFCPSessionModificationResponse) GetHeader() *PFCPHeader {
	return sr.Header
}
