package ie

import (
	"bytes"
	"testing"

	dt "github.com/fiorix/go-diameter/diam/datatype"
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

func TestCreatePDRStruct(t *testing.T) {
	ruleID := []byte{0x00, 0x10}
	pdrid := NewInformationElement(
		IEPDRID,
		0,
		dt.OctetString(ruleID),
	)

	pre := uint32(10)
	precedence := NewInformationElement(
		IEPrecedence,
		0,
		dt.Unsigned32(pre),
	)

	oHR := dt.OctetString([]byte{0x00})
	outerHeaderRemoval := NewInformationElement(
		IEOuterHeaderRemoval,
		0,
		oHR,
	)

	farId := dt.Unsigned32(100)
	farID := NewInformationElement(
		IEFARID,
		0,
		farId,
	)

	sourceinterface := uint8(1)
	si := NewInformationElement(
		IESourceInterface,
		0,
		dt.OctetString([]byte{sourceinterface}),
	)

	ft := []byte{0x01, 0x00, 0x00, 0xFF, 0xFF, 0xC0, 0xA8, 0x1, 0x65}
	fteid := NewInformationElement(
		IEFTEID,
		0,
		dt.OctetString(ft),
	)

	ni := dt.OctetString("internet.mnc012.mcc345.gprs")
	networkInstance := NewInformationElement(
		IENetworkInstance,
		0,
		ni,
	)
	pdi := NewPDI(&si, &fteid, &networkInstance, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	pdiB, err := pdi.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing PDI struct %+v", err)

	}

	pdiIE := NewInformationElement(
		IEPDI,
		0,
		dt.OctetString(pdiB),
	)

	createPDR := NewCreatePDR(&pdrid, &precedence, &pdiIE, &outerHeaderRemoval, &farID, nil, nil, nil)
	b, err := createPDR.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing CreatePDR %+v", err)

	}
	ba := []byte{0x00, 0x38, 0x00, 0x02, 0x00, 0x10,
		0x00, 0x1D, 0x00, 0x04, 0x00, 0x00, 0x00, 0xa,
		0x00, 0x02, 0x00, 0x31,
		0x00, 0x14, 0x00, 0x01, 0x01,
		0x00, 0x15, 0x00, 0x09, 0x01, 0x00, 0x00, 0xFF, 0xFF, 0xC0, 0xA8, 0x1, 0x65,
		0x00, 0x16, 0x00, 0x1b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x6d, 0x6e, 0x63, 0x30, 0x31, 0x32, 0x2e, 0x6d, 0x63, 0x63, 0x33, 0x34, 0x35, 0x2e, 0x67, 0x70, 0x72, 0x73,

		0x00, 0x5F, 0x00, 0x1, 0x00, 0x00, 0x6C, 0x00, 0x4, 0x00, 0x00, 0x00, 0x64,
	}

	if !bytes.Equal(b, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, b)

	}

}
