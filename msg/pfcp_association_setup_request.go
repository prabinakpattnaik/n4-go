package msg

import (
	"fmt"
	"net"
	"time"

	"bitbucket.org/sothy5/n4-go/ie"

	dt "github.com/fiorix/go-diameter/diam/datatype"
)

var (
	dataPlaneNodeID = []byte{0x0, 0xC0, 0xa8, 0x1, 0x20}
	ipv4address     = net.IPv4(8, 8, 8, 8)
)

//PFCPAssociationSetupRequest
type PFCPAssociationSetupRequest struct {
	Header                         *PFCPHeader
	NodeID                         *ie.InformationElement
	RecoveryTimeStamp              *ie.InformationElement
	UPFunctionFeatures             *ie.InformationElement
	CPFunctionFeatures             *ie.InformationElement
	UserPlaneIPResourceInformation ie.InformationElements
}

//NewPFCPAssociationSetupRequest creates new PFCPAssociationSetupRequst
func NewPFCPAssociationSetupRequest(h *PFCPHeader, n, r, u, c *ie.InformationElement, ui ie.InformationElements) PFCPAssociationSetupRequest {
	//if n.Type == ie.IEReserved || r.Type == ie.IEReserved {
	//	return nil
	//}
	return PFCPAssociationSetupRequest{

		Header:                         h,
		NodeID:                         n,
		RecoveryTimeStamp:              r,
		UPFunctionFeatures:             u,
		CPFunctionFeatures:             c,
		UserPlaneIPResourceInformation: ui,
	}

}

func FromPFCPMessage(m *PFCPMessage) (PFCP, error) {
	var n, r, u, c, cause ie.InformationElement
	var cpfseid, cPDR, cFAR, cURR, cQER, cBAR ie.InformationElement
	var offending, createdpdr ie.InformationElement
	var rPDR, rFAR, rURR, rQER, rBAR, rTE ie.InformationElement
	var uis ie.InformationElements

	for _, informationElement := range m.IEs {
		switch informationElement.Type {
		case ie.IENodeID:
			n = informationElement
		case ie.IERecoveryTimestamp:
			r = informationElement
		case ie.IEUPFunctionFeatures:
			u = informationElement
		case ie.IECPFunctionFeatures:
			c = informationElement
		case ie.IEUserPlaneIPResourceInformation:
			uis = append(uis, informationElement)
		case ie.IECause:
			cause = informationElement
		case ie.IEFSEID:
			cpfseid = informationElement
		case ie.IECreatePDR:
			cPDR = informationElement
		case ie.IECreateFAR:
			cFAR = informationElement
		case ie.IECreateURR:
			cURR = informationElement
		case ie.IECreateQER:
			cQER = informationElement
		case ie.IECreateBAR:
			cBAR = informationElement
		case ie.IEOffendingIE:
			offending = informationElement
		case ie.IECreatedPDR:
			createdpdr = informationElement
		case ie.IERemovePDR:
			rPDR = informationElement
		case ie.IERemoveFAR:
			rFAR = informationElement
		case ie.IERemoveURR:
			rURR = informationElement
		case ie.IERemoveQER:
			rQER = informationElement
		case ie.IERemoveBAR:
			rBAR = informationElement
		case ie.IERemoveTrafficEndpoint:
			rTE = informationElement

		default:
			return nil, fmt.Errorf("No matching needed Information Element")
		}

	}

	switch m.Header.MessageType {
	case AssociationSetupRequestType:
		pfcpAssociationSetupRequest := NewPFCPAssociationSetupRequest(m.Header, &n, &r, &u, &c, uis)
		return pfcpAssociationSetupRequest, nil
	case AssociationSetupResponseType:
		pfcpAssociationSetupResponse := NewPFCPAssociationSetupResponse(m.Header, &n, &cause, &r, &u, &c, uis)
		return pfcpAssociationSetupResponse, nil
	case SessionEstablishmentRequestType:
		pfcpSessionEstablishmentRequest := NewPFCPSessionEstablishmentRequest(m.Header, &n, &cpfseid, &cPDR, &cFAR, &cURR, &cQER, &cBAR, nil, nil, nil, nil, nil)
		return pfcpSessionEstablishmentRequest, nil
	case SessionEstablishmentResponseType:
		pfcpSessionEstablishmentResponse := NewPFCPSessionEstablishmentResponse(m.Header, &n, &cause, &offending, &cpfseid, &createdpdr, nil, nil, nil, nil)
		return pfcpSessionEstablishmentResponse, nil
	case SessionModificationRequestType:
		pfcpSessionModificationRequest := NewPFCPSessionModificationRequest(m.Header, &cpfseid, &rPDR, &rFAR, &rURR, &rQER, &rBAR, &rTE, &cPDR, &cFAR, &cURR, &cQER, &cBAR, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		return pfcpSessionModificationRequest, nil
	case SessionModificationResponseType:
		pfcpSessionModificationResponse := NewPFCPSessionModificationResponse(m.Header, &cause, &offending, &createdpdr)
		return pfcpSessionModificationResponse, nil
	case SessionDeletionResponseType:
		pfcpSessionDeletionResponse := NewPFCPSessionDeletionResponse(m.Header, &cause, &offending)
		return pfcpSessionDeletionResponse, nil
	default:
		return nil, fmt.Errorf("No matching PFCP Message Type")
	}

}

//ProcessAssociationSetupRequest process the PFCPMessage into PFCPAssociationSetupRequest, record the relvant details and create the right ProcessAssociationSetupResponse.
func ProcessAssociationSetupRequest(m *PFCPMessage) ([]byte, error) {
	//TODO: casting
	pfcpAssociationSetupRequest, err := FromPFCPMessage(m)
	if err != nil {
		return nil, err
	}
	// How to check all details are right.
	// record nodeID, timestamp
	// overlooked CPFunction Features.
	// PFCPAssociationSetupResponse Creation

	var length uint16

	n := ie.NewInformationElement(
		ie.IENodeID,
		0,
		dt.OctetString(dataPlaneNodeID),
	)
	length = n.Len() + ie.IEBasicHeaderSize
	c := ie.NewInformationElement(
		ie.IECause,
		0,
		dt.OctetString([]byte{0x01}),
	)
	length = length + c.Len() + ie.IEBasicHeaderSize

	r := ie.NewInformationElement(
		ie.IERecoveryTimestamp,
		0,
		dt.Time(time.Now()),
	)
	length = length + r.Len() + ie.IEBasicHeaderSize

	u := ie.NewUPFunctionFeatures(false, true, true, false, false, false, false, false, false, false, false, false, false, false, false)
	b, err := u.Serialize()
	if err != nil {
		return nil, err
	}

	upFunctionFeaturesIE := ie.NewInformationElement(
		ie.IEUPFunctionFeatures,
		0,
		dt.OctetString(b),
	)
	length = length + upFunctionFeaturesIE.Len() + ie.IEBasicHeaderSize
	upIPResourceInformation := ie.NewUPIPResourceInformation(true, false, 0, false, false, 0, ipv4address, nil, nil, 0)
	b, err = upIPResourceInformation.Serialize()
	if err != nil {
		return nil, err
	}

	ui := ie.NewInformationElement(
		ie.IEUserPlaneIPResourceInformation,
		0,
		dt.OctetString(b),
	)
	length = length + ui.Len() + ie.IEBasicHeaderSize

	length = length + PFCPBasicMessageSize

	header := pfcpAssociationSetupRequest.GetHeader()
	header.MessageType = AssociationSetupResponseType
	header.MessageLength = length
	var uis ie.InformationElements
	uis = append(uis, ui)
	pfcpAssociationSetupResponse := NewPFCPAssociationSetupResponse(header, &n, &c, &r, &upFunctionFeaturesIE, nil, uis)
	b, err = pfcpAssociationSetupResponse.Serialize()
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (ar PFCPAssociationSetupRequest) Serialize() ([]byte, error) {

	var b []byte
	if ar.NodeID == nil || ar.RecoveryTimeStamp == nil {
		return b, nil
	}

	output := make([]byte, ar.Len())
	pfcpend := uint16(PFCPBasicHeaderLength) + PFCPBasicMessageSize
	copy(output[:pfcpend], ar.Header.Serialize())
	nb, _ := ar.NodeID.Serialize()
	nodeIDEnd := pfcpend + ie.IEBasicHeaderSize + ar.NodeID.Len()
	copy(output[pfcpend:nodeIDEnd], nb)

	recoveryTimestampEnd := nodeIDEnd + ie.IEBasicHeaderSize + ar.RecoveryTimeStamp.Len()
	rb, _ := ar.RecoveryTimeStamp.Serialize()
	copy(output[nodeIDEnd:recoveryTimestampEnd], rb)

	var upFunctionFeaturesEnd, cpFunctionFeaturesEnd, upIPResourceInformationEnd uint16

	if ar.UPFunctionFeatures != nil {
		upFunctionFeaturesEnd = recoveryTimestampEnd + ie.IEBasicHeaderSize + ar.UPFunctionFeatures.Len()
		ub, _ := ar.UPFunctionFeatures.Serialize()
		copy(output[recoveryTimestampEnd:upFunctionFeaturesEnd], ub)
	}

	if ar.CPFunctionFeatures != nil {
		cb, _ := ar.CPFunctionFeatures.Serialize()
		if upFunctionFeaturesEnd == 0 {
			cpFunctionFeaturesEnd = recoveryTimestampEnd + ie.IEBasicHeaderSize + ar.CPFunctionFeatures.Len()
			copy(output[recoveryTimestampEnd:cpFunctionFeaturesEnd], cb)
		} else {
			cpFunctionFeaturesEnd = upFunctionFeaturesEnd + ie.IEBasicHeaderSize + ar.CPFunctionFeatures.Len()
			copy(output[upFunctionFeaturesEnd:cpFunctionFeaturesEnd], cb)
		}

	}
	if len(ar.UserPlaneIPResourceInformation) > 0 {
		for _, informationElement := range ar.UserPlaneIPResourceInformation {
			ib, err := informationElement.Serialize()
			if err != nil {
				return nil, err
			}

			if upFunctionFeaturesEnd > 0 && cpFunctionFeaturesEnd == 0 {
				upIPResourceInformationEnd = upFunctionFeaturesEnd + uint16(len(ib))
				copy(output[upFunctionFeaturesEnd:upIPResourceInformationEnd], ib)
				upFunctionFeaturesEnd = upIPResourceInformationEnd
			}
		}
	}

	return output, nil
}

func (ar PFCPAssociationSetupRequest) Len() uint16 {
	return uint16(PFCPBasicHeaderLength) + ar.Header.MessageLength
}

func (ar PFCPAssociationSetupRequest) Type() PFCPType {
	return ar.Header.MessageType
}

func (ar PFCPAssociationSetupRequest) GetHeader() *PFCPHeader {
	return ar.Header
}
