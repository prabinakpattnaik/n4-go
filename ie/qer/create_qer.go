package qer

import (
	"encoding/binary"
	"fmt"

	"github.com/prabinakpattnaik/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

//CreateQER is a struct of createQER
type CreateQER struct {
	QERID                 *ie.InformationElement
	QERCorrelationID      *ie.InformationElement
	GateStatus            *ie.InformationElement
	MBR                   *ie.InformationElement
	GBR                   *ie.InformationElement
	QFI                   *ie.InformationElement
	RQI                   *ie.InformationElement
	PagingPolicyIndicator *ie.InformationElement
	AverageWindow         *ie.InformationElement
}

func NewCreateQER(qerid uint32, qCID uint32, gs *GateStatus, m *BR, g *BR, qfi uint8, rqi bool) (*CreateQER, error) {
	if !(qerid > 0) {
		return nil, fmt.Errorf("Not valid QER ID")
	}
	qerIDIE := ie.NewInformationElement(
		ie.IEQERID,
		0,
		dt.Unsigned32(qerid),
	)

	var qCIDIE ie.InformationElement
	if qCID > 0 {
		qCIDIE = ie.NewInformationElement(
			ie.IEQERID,
			0,
			dt.Unsigned32(qCID),
		)
	}

	if gs == nil {
		return nil, fmt.Errorf("Mandatory field missing : GateStatus")
	}
	b, err := gs.Serialize()
	if err != nil {
		return nil, err
	}
	gsIE := ie.NewInformationElement(
		ie.IEGateStatus,
		0,
		dt.OctetString(b),
	)
	var mbrIE ie.InformationElement
	if m != nil {
		b, err := m.Serialize()
		if err != nil {
			return nil, err
		}

		mbrIE = ie.NewInformationElement(
			ie.IEMBR,
			0,
			dt.OctetString(b),
		)

	}

	var gbrIE ie.InformationElement
	if g != nil {
		b, err := g.Serialize()
		if err != nil {
			return nil, err
		}

		gbrIE = ie.NewInformationElement(
			ie.IEGBR,
			0,
			dt.OctetString(b),
		)

	}

	var gfiIE ie.InformationElement
	if (qfi > 0) && (qfi < 0x40) {
		gfiIE = ie.NewInformationElement(
			ie.IEQFI,
			0,
			dt.OctetString(qfi),
		)

	}

	var rqiIE ie.InformationElement
	if rqi {
		rqiIE = ie.NewInformationElement(
			ie.IERQI,
			0,
			dt.OctetString(qfi),
		)

	}

	createQER := &CreateQER{
		QERID:                 &qerIDIE,
		QERCorrelationID:      &qCIDIE,
		GateStatus:            &gsIE,
		MBR:                   &mbrIE,
		GBR:                   &gbrIE,
		QFI:                   &gfiIE,
		RQI:                   &rqiIE,
		PagingPolicyIndicator: nil,
		AverageWindow:         nil,
	}

	return createQER, nil

}

func (c CreateQER) Serialize() ([]byte, error) {
	if c.QERID == nil || c.QERID.Type == ie.IEReserved {
		return nil, fmt.Errorf("CreateQER does not have valid QERID")
	}
	b, err := c.QERID.Serialize()
	if err != nil {
		return nil, fmt.Errorf("CreateQER: QERID serialization error")
	}

	if c.QERCorrelationID != nil && c.QERCorrelationID.Type != ie.IEReserved {
		b1, err := c.QERCorrelationID.Serialize()
		if err != nil {
			return nil, fmt.Errorf("CreateQER: QERCorrelationID serialization error")
		}
		b = append(b, b1...)
	}

	if c.GateStatus == nil || c.GateStatus.Type == ie.IEReserved {
		return nil, fmt.Errorf("CreateQER does not have valid GateStatus")
	}
	b1, err := c.GateStatus.Serialize()
	if err != nil {
		return nil, fmt.Errorf("CreateQER: GateStatus serialization error")
	}
	b = append(b, b1...)

	if c.MBR != nil && c.MBR.Type != ie.IEReserved {
		b1, err := c.MBR.Serialize()
		if err != nil {
			return nil, fmt.Errorf("CreateQER: MBR serialization error")
		}
		b = append(b, b1...)
	}

	if c.GBR != nil && c.GBR.Type != ie.IEReserved {
		b1, err := c.GBR.Serialize()
		if err != nil {
			return nil, fmt.Errorf("CreateQER: GBR serialization error")
		}
		b = append(b, b1...)
	}

	if c.QFI != nil && c.QFI.Type != ie.IEReserved {
		b1, err := c.QFI.Serialize()
		if err != nil {
			return nil, fmt.Errorf("CreateQER: QFI serialization error")
		}
		b = append(b, b1...)
	}

	if c.RQI != nil && c.RQI.Type != ie.IEReserved {
		b1, err := c.RQI.Serialize()
		if err != nil {
			return nil, fmt.Errorf("CreateQER: GBR serialization error")
		}
		b = append(b, b1...)
	}

	return b, nil
}

type GateValue int

const (
	OPEN GateValue = iota
	CLOSED
)

type GateStatus struct {
	ULGate GateValue
	DLGate GateValue
}

func NewGateStatus(u, d GateValue) *GateStatus {
	return &GateStatus{
		ULGate: u,
		DLGate: d,
	}
}

func (g *GateStatus) Serialize() (byte, error) {
	var b uint8
	if g.DLGate == CLOSED {
		b = 1
	}
	if g.ULGate == CLOSED {
		b = b | 0x4
	}
	return b, nil
}

type BR struct {
	UL uint64
	DL uint64
}

func NewBR(u, d uint64) *BR {
	return &BR{
		UL: u,
		DL: d,
	}
}

func (b *BR) Serialize() ([]byte, error) {
	var r []byte
	d := make([]byte, 8)
	binary.BigEndian.PutUint64(d, b.UL)
	d1 := make([]byte, 8)
	binary.BigEndian.PutUint64(d1, b.DL)
	r = append(r, d[3:]...)
	r = append(r, d1[3:]...)
	return r, nil
}
