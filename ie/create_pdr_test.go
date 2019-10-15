package ie

import (
	"bytes"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	"testing"
)

func TestIECreatePDR(t *testing.T) {
	// included: PDR ID, Precedence, PDI, Outer Header Removal, FAR ID
	// not included: URRID, QER ID, Activate Pre-defined rules

	ruleID := []byte{0x00, 0x10}
	pdrid := NewInformationElement(
		IEPDRID,
		0,
		dt.OctetString(ruleID),
	)

	b1, err := pdrid.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	pre := uint32(10)

	precedence := NewInformationElement(
		IEPrecedence,
		0,
		dt.Unsigned32(pre),
	)

	b2, err := precedence.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	// PDI

	oHR := dt.OctetString([]byte{0x00})

	outerHeaderRemoval := NewInformationElement(
		IEOuterHeaderRemoval,
		0,
		oHR,
	)

	b3, err := outerHeaderRemoval.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	farId := dt.Unsigned32(100)

	farID := NewInformationElement(
		IEFARID,
		0,
		farId,
	)
	b4, err := farID.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	b1 = append(b1, b2...)
	b1 = append(b1, b3...)
	b1 = append(b1, b4...)

	i := NewInformationElement(

		IECreatePDR,
		0,
		dt.OctetString(b1),
	)

	if i.Length != uint16(len(b1)) {

		t.Fatalf("Unexpected length. want 2, have %d", i.Length)

	}

	ba := []byte{0x00, 0x01, 0x00, 0x1b, 0x00, 0x38, 0x00, 0x02, 0x00, 0x10,
		0x00, 0x1D, 0x00, 0x04, 0x00, 0x00, 0x00, 0xa, 0x00, 0x5F, 0x00, 0x1, 0x00, 0x00, 0x6C, 0x00, 0x4, 0x00, 0x00, 0x00, 0x64}

	b, err := i.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)

	}

	if !bytes.Equal(b, ba) {

		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, b)

	}

}
