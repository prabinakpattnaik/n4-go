package msg

import (
	"errors"

	"bitbucket.org/sothy5/n4-go/ie"
)

type PFCPAssociationSetupResponse struct {
	Header                         *PFCPHeader
	NodeID                         *ie.InformationElement
	Cause                          *ie.InformationElement
	RecoveryTimeStamp              *ie.InformationElement
	UPFunctionFeatures             *ie.InformationElement
	CPFunctionFeatures             *ie.InformationElement
	UserPlaneIPResourceInformation *ie.InformationElement
}

func NewPFCPAssociationSetupResponse(h *PFCPHeader, n, cause, r, u, c, ui *ie.InformationElement) PFCP {

	//if n == nil || r == nil {
	//	return nil
	//}
	return PFCPAssociationSetupResponse{
		Header:                         h,
		NodeID:                         n,
		Cause:                          cause,
		RecoveryTimeStamp:              r,
		UPFunctionFeatures:             u,
		CPFunctionFeatures:             c,
		UserPlaneIPResourceInformation: ui,
	}

}

func (res PFCPAssociationSetupResponse) Serialize() ([]byte, error) {
	if res.NodeID == nil || res.RecoveryTimeStamp == nil || res.Cause == nil {
		return nil, errors.New("Failed to serialize PFCPAssociation Setup Response:Either Nodeid, RecoveryTimestamp or Cause is nil ")
	}
	output := make([]byte, res.Len())
	pfcpend := uint16(PFCPBasicHeaderLength) + PFCPBasicMessageSize
	copy(output[:pfcpend], res.Header.Serialize())
	nb, err := res.NodeID.Serialize()
	if err != nil {
		return nil, err
	}
	nodeIDEnd := pfcpend + ie.IEBasicHeaderSize + res.NodeID.Len()
	copy(output[pfcpend:nodeIDEnd], nb)

	c, err := res.Cause.Serialize()
	if err != nil {
		return nil, err
	}

	causeEnd := nodeIDEnd + ie.IEBasicHeaderSize + res.Cause.Len()
	copy(output[nodeIDEnd:causeEnd], c)

	recoveryTimestampEnd := causeEnd + ie.IEBasicHeaderSize + res.RecoveryTimeStamp.Len()
	rb, err := res.RecoveryTimeStamp.Serialize()
	if err != nil {
		return nil, err
	}
	copy(output[causeEnd:recoveryTimestampEnd], rb)

	var upFunctionFeaturesEnd, cpFunctionFeaturesEnd uint16

	if res.UPFunctionFeatures != nil {
		upFunctionFeaturesEnd = recoveryTimestampEnd + ie.IEBasicHeaderSize + res.UPFunctionFeatures.Len()
		ub, err := res.UPFunctionFeatures.Serialize()
		if err != nil {
			return nil, err
		}
		copy(output[recoveryTimestampEnd:upFunctionFeaturesEnd], ub)

	}

	if res.CPFunctionFeatures != nil {
		cb, err := res.CPFunctionFeatures.Serialize()

		if err != nil {
			return nil, err
		}
		if upFunctionFeaturesEnd == 0 {

			cpFunctionFeaturesEnd = recoveryTimestampEnd + ie.IEBasicHeaderSize + res.CPFunctionFeatures.Len()
			copy(output[recoveryTimestampEnd:cpFunctionFeaturesEnd], cb)
		} else {
			cpFunctionFeaturesEnd = upFunctionFeaturesEnd + ie.IEBasicHeaderSize + res.CPFunctionFeatures.Len()
			copy(output[upFunctionFeaturesEnd:cpFunctionFeaturesEnd], cb)
		}

	}
	if res.UserPlaneIPResourceInformation != nil {
		ib, err := res.UserPlaneIPResourceInformation.Serialize()

		if err != nil {
			return nil, err
		}
		if upFunctionFeaturesEnd > 0 && cpFunctionFeaturesEnd == 0 {
			//upIPResourceInformationEnd = upFunctionFeaturesEnd + ie.IEBasicHeaderSize + res.UserPlaneIPResourceInformation.Len()
			copy(output[upFunctionFeaturesEnd:], ib)
		}
	}

	return output, nil
}

func (res PFCPAssociationSetupResponse) Len() uint16 {
	return uint16(PFCPBasicHeaderLength) + res.Header.MessageLength

}

func (res PFCPAssociationSetupResponse) Type() PFCPType {
	return res.Header.MessageType
}
func (ar PFCPAssociationSetupResponse) GetHeader() *PFCPHeader {
	return ar.Header
}
