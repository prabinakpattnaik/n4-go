package ie

import (
	"fmt"
)

// CreateFAR struct
type CreateFARWithIE struct {
	FARID                 *InformationElement
	ApplyAction           *InformationElement
	ForwardingParameters  *InformationElement
	DuplicatingParameters *InformationElement
	BARID                 *InformationElement
}

//NewCreateFAR creates new CreateFAR struct
func NewCreateFAR(farid, applyAction, forwardingParameters, duplicatingParameters, barid *InformationElement) *CreateFARWithIE {
	return &CreateFARWithIE{
		FARID:                 farid,
		ApplyAction:           applyAction,
		ForwardingParameters:  forwardingParameters,
		DuplicatingParameters: duplicatingParameters,
		BARID:                 barid,
	}

}

//Serialize function convert CreateFAR struct into byte array
func (c CreateFARWithIE) Serialize() ([]byte, error) {
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

	if c.ForwardingParameters != nil && c.ForwardingParameters.Type != IEReserved {
		b1, err = c.ForwardingParameters.Serialize()
		if err != nil {
			return nil, fmt.Errorf("Forwarding parameters serialization error")
		}
		b = append(b, b1...)
	}

	if c.DuplicatingParameters != nil && c.DuplicatingParameters.Type != IEReserved {
		b1, err = c.DuplicatingParameters.Serialize()
		if err != nil {
			return nil, fmt.Errorf("DuplicatingParameters serialization error")
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
