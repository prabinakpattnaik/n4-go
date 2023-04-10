package ie

import (
	"encoding/binary"
	"fmt"

	dt "github.com/fiorix/go-diameter/diam/datatype"
)

type UpdatePDR struct {
	PDRID      *PacketDetectionRuleID `tlv:"56"`
	OHR        *OuterHeaderRemoval    `tlv:"95"`
	Precedence *Precedence            `tlv:"29"`
	Pdi        *PDI                   `tlv:"2"`
	FarID      *FARID                 `tlv:"108"`
	UrrID      *uint32                `tlv:"81"`
	//QERID                   *QERID                   `tlv:"109"`
	//ActivatePredefinedRules *ActivatePredefinedRules `tlv:"106"`
	//DeactivatePredefinedRules *DeactivatePredefinedRules `tlv:"107"`

}

//UpdatePDR produces binary for IEUpdatePDR
func (c *UpdatePDR) Serialize() ([]byte, error) {
	if c.PDRID == nil {
		return nil, fmt.Errorf("No valid PDRID")
	}
	d := make([]byte, 2)
	binary.BigEndian.PutUint16(d, c.PDRID.RuleId)
	pdrid := NewInformationElement(
		IEPDRID,
		0,
		dt.OctetString(d),
	)

	var ohr InformationElement
        if c.OHR != nil {
        ohr = NewInformationElement(
                IEOuterHeaderRemoval,
                0,
                dt.OctetString([]byte{c.OHR.OuterHeaderRemovalDescription}),
              )
        }

	if c.Precedence == nil {
		return nil, fmt.Errorf("No valid Precedence")
	}
	precedence := NewInformationElement(
		IEPrecedence,
		0,
		dt.Unsigned32(c.Precedence.PrecedenceValue),
	)
	if c.Pdi == nil {
		return nil, fmt.Errorf("No valid PDI")
	}

	b, err := c.Pdi.Serialize()
	var pdi InformationElement
	if err == nil {
		pdi = NewInformationElement(
			IEPDI,
			0,
			dt.OctetString(b),
		)
	}
	farid := NewInformationElement(
		IEFARID,
		0,
		dt.Unsigned32(c.FarID.FarIdValue),
	)

	var urrid InformationElement
	if c.UrrID != nil {
		urrid = NewInformationElement(
			IEURRID,
			0,
			dt.Unsigned32(uint32(*c.UrrID)),
		)
	}
	updatePDR := NewUpdatePDR(&pdrid, &ohr, &precedence, &pdi, &farid, &urrid, nil, nil, nil)
	return updatePDR.Serialize()

}
