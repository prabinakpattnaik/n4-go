package msg

import (
	"fmt"

	"bitbucket.org/sothy5/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

type PFCPSessionDeletionResponse struct {
	Header                     *PFCPHeader
	Cause                      *ie.InformationElement
	OffendingIE                *ie.InformationElement
	LoadControlInformation     *ie.InformationElement
	OverloadControlInformation *ie.InformationElement
	UsageReport                *ie.InformationElement
}

func NewPFCPSessionDeletionResponse(h *PFCPHeader, c, offending *ie.InformationElement) PFCPSessionDeletionResponse {
	return PFCPSessionDeletionResponse{
		Header:      h,
		Cause:       c,
		OffendingIE: offending,
	}
}

//Serialize function converts PFCPSessionDeletionResponse into byte array
func (smr PFCPSessionDeletionResponse) Serialize() ([]byte, error) {
	var b []byte
	if smr.Cause == nil {
		return nil, fmt.Errorf("There is no Cause value")
	}

	dataLength := smr.Len()
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

	copy(output[pfcpend:], b)
	//TODO: remaining to be added (load control information,OverloadControlInformation, UsageReport,FailedRuleID, and CreatedTrafficEndpoint )

	return output, nil

}

func ProcessPFCPSessionDeletionRequest(m *PFCPMessage, sSEID uint64) ([]byte, error) {
	if m.Header.MessageType != SessionDeletionRequestType {
		return nil, fmt.Errorf("PFCP Session Deletion Request could not found")
	}

	c := ie.NewInformationElement(
		ie.IECause,
		0,
		dt.OctetString([]byte{0x01}),
	)
	length := ie.IEBasicHeaderSize + c.Len()
	pfcpHeader := NewPFCPHeader(1, false, true, SessionDeletionResponseType, length+12, sSEID, m.Header.SequenceNumber, 0)
	pfcpSessionDeletionResponse := NewPFCPSessionDeletionResponse(pfcpHeader, &c, nil)
	b, err := pfcpSessionDeletionResponse.Serialize()
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (sr PFCPSessionDeletionResponse) Len() uint16 {
	return uint16(PFCPBasicHeaderLength) + sr.Header.MessageLength
}

func (sr PFCPSessionDeletionResponse) Type() PFCPType {
	return sr.Header.MessageType
}

func (sr PFCPSessionDeletionResponse) GetHeader() *PFCPHeader {
	return sr.Header
}
