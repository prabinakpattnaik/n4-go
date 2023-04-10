package ie

import (
	"fmt"
)

// UpdateFAR struct
type UpdateFARWithIE struct {
	FARID                       *InformationElement
	ApplyAction                 *InformationElement
	UpdateForwardingParameters  *InformationElement
	UpdateDuplicatingParameters *InformationElement
	BARID                       *InformationElement
}

//NewUpdateFAR creates new UpdateFAR struct
func NewUpdateFAR(farid, applyAction, updateforwardingParameters, updateduplicatingParameters, barid *InformationElement) *UpdateFARWithIE {
	return &UpdateFARWithIE{
		FARID:                       farid,
		ApplyAction:                 applyAction,
		UpdateForwardingParameters:  updateforwardingParameters,
		UpdateDuplicatingParameters: updateduplicatingParameters,
		BARID:                       barid,
	}

}

//Serialize function convert UpdateFAR struct into byte array
func (c UpdateFARWithIE) Serialize() ([]byte, error) {
	if c.FARID == nil || c.FARID.Type == IEReserved {
		return nil, fmt.Errorf("CreateFAR does not have valid FARID")
	}
	b, err := c.FARID.Serialize()
	if err != nil {
		return nil, fmt.Errorf("FARID serialization error")
	}

	if c.ApplyAction == nil || c.ApplyAction.Type == IEReserved {
		return nil, fmt.Errorf("CreateFAR does not have valid ApplyAction")
	}
	b1, err := c.ApplyAction.Serialize()
	if err != nil {
		return nil, fmt.Errorf("ApplyAction serialization error")
	}
	b = append(b, b1...)

	if c.UpdateForwardingParameters != nil && c.UpdateForwardingParameters.Type != IEReserved {
		b1, err = c.UpdateForwardingParameters.Serialize()
		if err != nil {
			return nil, fmt.Errorf("Update forwarding parameters serialization error")
		}
		b = append(b, b1...)
	}

	if c.UpdateDuplicatingParameters != nil && c.UpdateDuplicatingParameters.Type != IEReserved {
		b1, err = c.UpdateDuplicatingParameters.Serialize()
		if err != nil {
			return nil, fmt.Errorf("Update duplicatingParameters serialization error")
		}
		b = append(b, b1...)
	}

	if c.BARID != nil && c.BARID.Type != IEReserved {
		b1, err = c.BARID.Serialize()
		if err != nil {
			return nil, fmt.Errorf("BARID serialization error")
		}
		b = append(b, b1...)
	}
	return b, nil
}
