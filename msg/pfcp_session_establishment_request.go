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
	CreatePDR                *ie.InformationElements
	CreateFAR                *ie.InformationElements
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
func NewPFCPSessionEstablishmentRequest(h *PFCPHeader, n, cpfseid *ie.InformationElement, createPDR, createFAR *ie.InformationElements, createURR, createQER, createBAR, createTrafficEndpoint, pdnType, userPlaneInactivityTimer, userID, traceInformation *ie.InformationElement) PFCPSessionEstablishmentRequest {
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
	pfcpend := uint16(PFCPBasicHeaderLength) + PFCPMessageSize
	output := make([]byte, pfcpend)
	copy(output[:pfcpend], sr.Header.Serialize())

	nb, err := sr.NodeID.Serialize()
	if err != nil {
		return nil, err
	}
	output = append(output, nb...)

	seid, err := sr.CPFSEID.Serialize()
	if err != nil {
		return nil, err
	}
	output = append(output, seid...)

	for _, pdrie := range *sr.CreatePDR {
		createpdr, err := pdrie.Serialize()
		if err != nil {
			return nil, err
		}
		output = append(output, createpdr...)
	}

	for _, farie := range *sr.CreateFAR {
		createfar, err := farie.Serialize()
		if err != nil {
			return nil, err
		}
		output = append(output, createfar...)
	}

	if sr.CreateURR != nil && sr.CreateURR.Type !=ie.IEReserved {
		createURR, err := sr.CreateURR.Serialize()
		if err != nil {
			return nil, err
		}
		output = append(output, createURR...)
	}

	if sr.CreateQER != nil {
		createQER, _ := sr.CreateQER.Serialize()
		if err != nil {
			return nil, err
		}
		output = append(output, createQER...)
	}

	if sr.UserPlaneInactivityTimer != nil && sr.UserPlaneInactivityTimer.Type != ie.IEReserved {
		upInactivityTimer, err := sr.UserPlaneInactivityTimer.Serialize()
		if err != nil {
			return nil, err
		}
		output = append(output, upInactivityTimer...)
	}
	if sr.PDNType != nil && sr.PDNType.Type != ie.IEReserved {
		pdnType, err := sr.PDNType.Serialize()
		if err != nil {
			return nil, err
		}
		output = append(output, pdnType...)
	}
	/*if int(dataLength)+PFCPBasicHeaderLength != len(output) {
	  return nil, fmt.Errorf("Length is wrong in the message")
	}
	*/
	return output, nil
	//TODO: remaining to be added

}

func ProcessPFCPSessionEstablishmentRequest(m *PFCPMessage, nodeIP net.IP, seid uint64) ([]byte, error) {
	pfcpMessage, err := FromPFCPMessage(m)
	if err != nil {
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

	var cPDR *ie.CreatePDRWithIE    
	for _,ieCreatePDR :=range * pfcpSessionEstablishmentRequest.CreatePDR {
	cPDR, err = ie.CreatePDRIEsFromBytes(ieCreatePDR.Data.Serialize())
	if err != nil {
		return nil, err
	}
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
