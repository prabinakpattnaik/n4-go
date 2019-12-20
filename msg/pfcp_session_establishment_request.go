package msg

import (
	"fmt"
	"net"

	"bitbucket.org/sothy5/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
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
	fmt.Printf("dataLength %d\n", dataLength)

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

	if sr.CreateURR != nil {
		createURR, _ := sr.CreateURR.Serialize()
		createUrrEnd := createfarEnd + ie.IEBasicHeaderSize + sr.CreateURR.Len()
		copy(output[createfarEnd:createUrrEnd], createURR)

	}

	return output, nil

	//TODO: remaining to be added

}

func ProcessPFCPSessionEstablishmentRequest(m *PFCPMessage, nodeIP net.IP, seid uint64) ([]byte, error) {
	pfcpMessage, err := FromPFCPMessage(m)
	if err != nil {
		fmt.Printf("pfcpSessionER:Header %+v\n", pfcpMessage)
		return nil, err
	}
	pfcpSessionEstablishmentRequest, ok := pfcpMessage.(PFCPSessionEstablishmentRequest)

	if !ok {
		return nil, fmt.Errorf("PFCP Session Establishment Request could not type asserted")
	}

	//TODO, add more business logic to check the details,CreatePDR, CreateFAR

	nodeID := []byte{0x00}
	nodeID = append(nodeID, nodeIP.To4()...)

	n := ie.NewInformationElement(
		ie.IENodeID,
		0,
		dt.OctetString(nodeID),
	)
	length := ie.IEBasicHeaderSize + n.Len()

	c := ie.NewInformationElement(
		ie.IECause,
		0,
		dt.OctetString([]byte{0x01}),
	)
	length = length + ie.IEBasicHeaderSize + c.Len()

	fseid := ie.NewFSEID(true, false, seid, nodeIP, nil)
	bb, err := fseid.Serialize()
	if err != nil {
		return nil, err
	}
	upfseid := ie.NewInformationElement(
		ie.IEFSEID,
		0,
		dt.OctetString(bb),
	)
	length = length + ie.IEBasicHeaderSize + upfseid.Len()

	cPDR, err := ie.CreatePDRIEsFromBytes(pfcpSessionEstablishmentRequest.CreatePDR.Data.Serialize())
	if err != nil {
		return nil, err
	}

	b, err := cPDR.PDRID.Serialize()
	if err != nil {
		return nil, err
	}
	createdPDRIE := ie.NewInformationElement(
		ie.IECreatedPDR,
		0,
		dt.OctetString(b),
	)
	length = length + ie.IEBasicHeaderSize + createdPDRIE.Len()

	pfcpHeader := NewPFCPHeader(1, false, true, SessionEstablishmentResponseType, length+12, pfcpSessionEstablishmentRequest.Header.SessionEndpointIdentifier, pfcpSessionEstablishmentRequest.Header.SequenceNumber, 0)
	pfcpSessionEstablishmentResponse := NewPFCPSessionEstablishmentResponse(pfcpHeader, &n, &c, nil, &upfseid, &createdPDRIE, nil, nil, nil, nil)
	b, err = pfcpSessionEstablishmentResponse.Serialize()
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
