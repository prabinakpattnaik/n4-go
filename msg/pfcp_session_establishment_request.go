package msg

import (
	"fmt"

	"bitbucket.org/sothy5/n4-go/ie"
)

//PFCPSessionEstablishmentRequest
type PFCPSessionEstablishmentRequest struct {
	Header                   *PFCPHeader
	NodeID                   *ie.InformationElement
	CPFSEID                  *ie.InformationElement
	CreatePDR                *ie.InformationElement
	CreateFAR                *ie.InformationElement
	CreateURR                *ie.InformationElement
	CreateQER                *ie.InformationElement
	CreateBAR                *ie.InformationElement
	CreateTrafficEndpoint    *ie.InformationElement
	PDNType                  *ie.InformationElement
	UserPlaneInactivityTimer *ie.InformationElement
	UserID                   *ie.InformationElement
	TraceInformation         *ie.InformationElement
}

//NewPFCPSessionEstablishmentRequest creates new PFCPSessionEstablishmentRequst
func NewPFCPSessionEstablishmentRequest(h *PFCPHeader, n, cpfseid, createPDR, createFAR, createURR, createQER, createBAR, createTrafficEndpoint, pdnType, userPlaneInactivityTimer, userID, traceInformation *ie.InformationElement) PFCPSessionEstablishmentRequest {
	//if n.Type == ie.IEReserved || r.Type == ie.IEReserved {
	//	return nil
	//}
	return PFCPSessionEstablishmentRequest{

		Header:                   h,
		NodeID:                   n,
		CPFSEID:                  cpfseid,
		CreatePDR:                createPDR,
		CreateFAR:                createFAR,
		CreateURR:                createURR,
		CreateQER:                createQER,
		CreateBAR:                createBAR,
		CreateTrafficEndpoint:    createTrafficEndpoint,
		PDNType:                  pdnType,
		UserPlaneInactivityTimer: userPlaneInactivityTimer,
		UserID:                   userID,
		TraceInformation:         traceInformation,
	}

}

func (sr PFCPSessionEstablishmentRequest) Serialize() ([]byte, error) {
	var b []byte
	if sr.NodeID == nil || sr.CPFSEID == nil || sr.CreatePDR == nil || sr.CreateFAR == nil {
		return b, nil
	}

	dataLength := sr.Len()

	output := make([]byte, dataLength)
	pfcpend := uint16(PFCPBasicHeaderLength) + PFCPMessageSize
	copy(output[:pfcpend], sr.Header.Serialize())

	nb, _ := sr.NodeID.Serialize()
	nodeIDEnd := pfcpend + ie.IEBasicHeaderSize + sr.NodeID.Len()
	copy(output[pfcpend:nodeIDEnd], nb)

	seid, _ := sr.CPFSEID.Serialize()
	seidEnd := nodeIDEnd + ie.IEBasicHeaderSize + sr.CPFSEID.Len()
	copy(output[nodeIDEnd:seidEnd], seid)

	createpdr, _ := sr.CreatePDR.Serialize()
	createpdrEnd := seidEnd + ie.IEBasicHeaderSize + sr.CreatePDR.Len()
	copy(output[seidEnd:createpdrEnd], createpdr)

	createfar, _ := sr.CreateFAR.Serialize()
	createfarEnd := createpdrEnd + ie.IEBasicHeaderSize + sr.CreateFAR.Len()
	copy(output[createpdrEnd:createfarEnd], createfar)

	return output, nil

	//TODO: remaining to be added

}

func ProcessPFCPSessionEstablishmentRequest(m *PFCPMessage) ([]byte, error) {
	pfcpMessage, err := FromPFCPMessage(m)
	if err != nil {
		return nil, err
	}
	pfcpSessionEstablishmentRequest, ok := pfcpMessage.(PFCPSessionEstablishmentRequest)
	if !ok {
		return nil, fmt.Errorf("PFCP Session Establishment Request could not type asserted")
	}
	b, err := pfcpSessionEstablishmentRequest.Serialize()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (sr PFCPSessionEstablishmentRequest) Len() uint16 {
	return uint16(PFCPBasicHeaderLength) + sr.Header.MessageLength
}

func (sr PFCPSessionEstablishmentRequest) Type() PFCPType {
	return sr.Header.MessageType
}

func (sr PFCPSessionEstablishmentRequest) GetHeader() *PFCPHeader {
	return sr.Header
}
