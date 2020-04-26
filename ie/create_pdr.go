package ie

import (
	"encoding/binary"
	"fmt"

	dt "github.com/fiorix/go-diameter/diam/datatype"
)

type CreatePDR struct {
	PDRID      *PacketDetectionRuleID `tlv:"56"`
	Precedence *Precedence            `tlv:"29"`
	Pdi        *PDI                   `tlv:"2"`
	OHR        *OuterHeaderRemoval    `tlv:"95"`
	FarID      *FARID                 `tlv:"108"`
	UrrID      *uint32                `tlv:"81"`
	//QERID                   *QERID                   `tlv:"109"`
	//ActivatePredefinedRules *ActivatePredefinedRules `tlv:"106"`

}

//CreatePDR produces binary for IECreatePDR
func (c *CreatePDR) Serialize() ([]byte, error) {
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
	var ohr InformationElement
	if c.OHR != nil {
	ohr = NewInformationElement(
		IEOuterHeaderRemoval,
		0,
		dt.OctetString([]byte{c.OHR.OuterHeaderRemovalDescription}),
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
	createPDR := NewCreatePDR(&pdrid, &precedence, &pdi, &ohr, &farid, &urrid, nil, nil)
	return createPDR.Serialize()

}
