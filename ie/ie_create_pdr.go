package ie

import (
	"fmt"
)

type CreatePDRWithIE struct {
	PDRID                   *InformationElement
	Precedence              *InformationElement
	PDI                     *InformationElement
	OuterHeaderRemoval      *InformationElement
	FARID                   *InformationElement
	URRID                   *InformationElement
	QERID                   *InformationElement
	ActivatePredefinedRules *InformationElement
}

func NewCreatePDR(pdrid, precedence, pdi, outerHeaderRemoval, farid, urrid, qerid, activatepredefinedRules *InformationElement) *CreatePDRWithIE {
	return &CreatePDRWithIE{
		PDRID:                   pdrid,
		Precedence:              precedence,
		PDI:                     pdi,
		OuterHeaderRemoval:      outerHeaderRemoval,
		FARID:                   farid,
		URRID:                   urrid,
		QERID:                   qerid,
		ActivatePredefinedRules: activatepredefinedRules,
	}

}

func (c CreatePDRWithIE) Serialize() ([]byte, error) {

	if c.PDRID == nil || c.PDRID.Type == IEReserved {
		return nil, fmt.Errorf("CreatePDR does not have valid PDRID")
	}
	b, err := c.PDRID.Serialize()
	if err != nil {
		return nil, fmt.Errorf("PDRID serialization error")
	}

	if c.Precedence == nil || c.Precedence.Type == IEReserved {
		return nil, fmt.Errorf("CreatePDR does not have valid Precedence")
	}
	b1, err := c.Precedence.Serialize()
	if err != nil {
		return nil, fmt.Errorf("Precedence serialization error")
	}
	b = append(b, b1...)

	if c.PDI == nil || c.PDI.Type == IEReserved {
		return nil, fmt.Errorf("CreatePDR does not have valid PDI")
	}
	b1, err = c.PDI.Serialize()
	if err != nil {
		return nil, fmt.Errorf("PDI serialization error")
	}
	b = append(b, b1...)

	if c.OuterHeaderRemoval != nil && c.OuterHeaderRemoval.Type != IEReserved {
		b1, err = c.OuterHeaderRemoval.Serialize()
		if err != nil {
			return nil, fmt.Errorf("OuterHederRemoval serialization error")
		}
		b = append(b, b1...)
	}

	if c.FARID != nil && c.FARID.Type != IEReserved {
		b1, err = c.FARID.Serialize()
		if err != nil {
			return nil, fmt.Errorf("FARID serialization error")
		}
		b = append(b, b1...)
	}
	if c.URRID != nil && c.URRID.Type != IEReserved {
		b1, err = c.URRID.Serialize()
		if err != nil {
			return nil, fmt.Errorf("URRID serialization error")
		}
		b = append(b, b1...)
	}

	if c.QERID != nil && c.QERID.Type != IEReserved {
		b1, err = c.QERID.Serialize()
		if err != nil {
			return nil, fmt.Errorf("QER ID serialization error")
		}
		b = append(b, b1...)
	}
	//TODO check correlation FARID and Activate PredefinedRules
	if c.ActivatePredefinedRules != nil && c.ActivatePredefinedRules.Type != IEReserved {
		b1, err = c.ActivatePredefinedRules.Serialize()
		if err != nil {
			return nil, fmt.Errorf("ActivatePredefinedRules serialization error")
		}
		b = append(b, b1...)
	}

	return b, nil

}

func CreatePDRIEsFromBytes(b []byte) (*CreatePDRWithIE, error) {
	var createPDRIEs InformationElements
	var pdrid, precedence, pdi, ohr, farid, urrid, qerid, activatePredefinedRules *InformationElement
	err := createPDRIEs.FromBytes(b)
	if err != nil {
		return nil, err
	}

	for _, informationElement := range createPDRIEs {
		switch informationElement.Type {
		case IEPDRID:
			pdrid = &informationElement
		case IEPrecedence:
			precedence = &informationElement
		case IEPDI:
			pdi = &informationElement
		case IEOuterHeaderRemoval:
			ohr = &informationElement
		case IEFARID:
			farid = &informationElement
		case IEURRID:
			urrid = &informationElement
		case IEQERID:
			qerid = &informationElement
		case IEActivatePredefinedRules:
			activatePredefinedRules = &informationElement

		default:
			return nil, fmt.Errorf("No matching needed Information Element for createPDR")
		}

	}
	cPDR := NewCreatePDR(pdrid, precedence, pdi, ohr, farid, urrid, qerid, activatePredefinedRules)
	return cPDR, nil

}
