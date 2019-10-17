package ie

import (
	"bytes"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	"testing"
)

func TestIECreateFAR(t *testing.T) {
	// included: FARID,ApplyAction,
	// not included: Forwarding Parameters, Duplicating Parameters, BAR ID

	farID := dt.Unsigned32(100)

	f := NewInformationElement(
		IEFARID,
		0,
		farID,
	)

	b1, err := f.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	applyAction := dt.OctetString([]byte{0x02})
	a := NewInformationElement(
		IEApplyAction,
		0,
		applyAction,
	)

	b2, err := a.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v\n", err)
	}

	b1 = append(b1, b2...)

	cf := NewInformationElement(
		IECreateFAR,
		0,
		dt.OctetString(b1),
	)

	if cf.Length != uint16(len(b1)) {
		t.Fatalf("Unexpected length. want %d, have %d", len(b1), cf.Length)
	}

	ba := []byte{0x00, 0x3, 0x00, 0x0d,
		0x00, 0x6C, 0x00, 0x4, 0x00, 0x00, 0x00, 0x64,
		0x00, 0x2C, 0x00, 0x01, 0x02,
	}

	bb, err := cf.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)

	}

	if !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

}
