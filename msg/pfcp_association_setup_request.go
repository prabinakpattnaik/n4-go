package msg

import (
	"bitbucket.org/sothy5/n4-go/ie"
)

type PFCPAssociationSetupRequest struct {
	Header                         *PFCPHeader
	NodeID                         *ie.InformationElement
	RecoveryTimeStamp              *ie.InformationElement
	UPFunctionFeatures             *ie.InformationElement
	CPFunctionFeatures             *ie.InformationElement
	UserPlaneIPResourceInformation *ie.InformationElement
}

func NewPFCPAssociationSetupRequest(h *PFCPHeader, n, r, u, c, ui *ie.InformationElement) *PFCPAssociationSetupRequest {

	if n == nil || r == nil {
		return nil
	}
	return &PFCPAssociationSetupRequest{
		Header:                         h,
		NodeID:                         n,
		RecoveryTimeStamp:              r,
		UPFunctionFeatures:             u,
		CPFunctionFeatures:             c,
		UserPlaneIPResourceInformation: ui,
	}

}

func (ar PFCPAssociationSetupRequest) Serialize() []byte {
	var b []byte
	if ar.NodeID == nil || ar.RecoveryTimeStamp == nil {
		return b
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
	if ar.UserPlaneIPResourceInformation != nil {
		ib, _ := ar.UserPlaneIPResourceInformation.Serialize()

		if upFunctionFeaturesEnd > 0 && cpFunctionFeaturesEnd == 0 {
			upIPResourceInformationEnd = upFunctionFeaturesEnd + ie.IEBasicHeaderSize + ar.UserPlaneIPResourceInformation.Len()
			copy(output[upFunctionFeaturesEnd:upIPResourceInformationEnd], ib)
		}
	}

	return output
}

func (ar PFCPAssociationSetupRequest) Len() uint16 {
	return uint16(PFCPBasicHeaderLength) + ar.Header.MessageLength
}
