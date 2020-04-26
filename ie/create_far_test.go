package ie

import (
	"bytes"
	"net"
	"testing"

	dt "github.com/fiorix/go-diameter/diam/datatype"
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

func TestCreateFARStruct(t *testing.T) {
	farID := dt.Unsigned32(100)

	f := NewInformationElement(
		IEFARID,
		0,
		farID,
	)

	applyAction := dt.OctetString([]byte{0x02})
	a := NewInformationElement(
		IEApplyAction,
		0,
		applyAction,
	)

	destinationinterface := uint8(1)
	d := NewInformationElement(
		IEDestinationInterface,
		0,
		dt.OctetString([]byte{destinationinterface}),
	)

	ohcd := uint8(1)
	teid := uint32(100)
	ip4address := net.ParseIP("8.8.8.8")
	ohc := NewOuterHeaderCreation(ohcd, teid, ip4address, nil, 0)
	b, err := ohc.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	ohcIE := NewInformationElement(
		IEOuterHeaderCreation,
		0,
		dt.OctetString(b),
	)

	fp := NewForwardingParameters(&d, nil, nil, &ohcIE, nil, nil, nil, nil, nil)
	bb, err := fp.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	fpIE := NewInformationElement(
		IEForwardingParameters,
		0,
		dt.OctetString(bb),
	)
	createFAR := NewCreateFAR(&f, &a, &fpIE, nil, nil)

	ba := []byte{0x00, 0x6C, 0x00, 0x4, 0x00, 0x00, 0x00, 0x64,
		0x00, 0x2C, 0x00, 0x01, 0x02,
		0x00, 0x04, 0x00, 0x12, 0x00, 0x2a, 0x00, 0x01, 0x01, 0x00, 0x54, 0x00, 0x09, 0x01, 0x00, 0x00, 0x00, 0x64, 0x08, 0x08, 0x08, 0x08}

	bb, err = createFAR.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)

	}

	if !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

}

func TestCreateFARNewStruct(t *testing.T) {
	farId := FARID{
		FarIdValue: 01,
	}
	aa := ApplyAction{
		Forw: true,
	}
	createFAR := CreateFAR{
		FarID:       &farId,
		ApplyAction: &aa,
	}
	b, err := createFAR.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}
	ba := []byte{0x00, 0x6c, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x2C, 0x00, 0x01, 0x02,
	}
	if !bytes.Equal(b, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, b)
	}

}
