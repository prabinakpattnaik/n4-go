package ie

import (
	"fmt"
)

type CreatedPDR struct {
	PDRID      *InformationElement
	LocalFTEID *InformationElement
}

func NewCreatedPDR(pdrid, localFTEID *InformationElement) *CreatedPDR {
	return &CreatedPDR{
		PDRID:      pdrid,
		LocalFTEID: localFTEID,
	}
}
func (c CreatedPDR) Serialize() ([]byte, error) {
	if c.PDRID == nil || c.PDRID.Type == IEReserved {
		return nil, fmt.Errorf("CreatedPDR does not have valid PDRID")
	}
	b, err := c.PDRID.Serialize()
	if err != nil {
		return nil, fmt.Errorf("PDRID serialization error")
	}

	if c.LocalFTEID != nil && c.LocalFTEID.Type != IEReserved {
		b1, err := c.LocalFTEID.Serialize()
		if err != nil {
			return nil, fmt.Errorf("localFTEID serialization error")
		}
		b = append(b, b1...)
	}
	return b, nil

}
