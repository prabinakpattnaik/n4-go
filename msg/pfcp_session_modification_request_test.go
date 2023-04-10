package msg

import (
	"bytes"
	"net"
	"testing"

	"github.com/prabinakpattnaik/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

func TestPFCPSessionModificationRequest(t *testing.T) {
	sei := uint64(400)
	sn := uint32(200)

	ruleID := []byte{0x00, 0x10}
	pdrid := ie.NewInformationElement(
		ie.IEPDRID,
		0,
		dt.OctetString(ruleID),
	)

	pre := uint32(10)
	precedence := ie.NewInformationElement(
		ie.IEPrecedence,
		0,
		dt.Unsigned32(pre),
	)

	farId := dt.Unsigned32(100)
	farID := ie.NewInformationElement(
		ie.IEFARID,
		0,
		farId,
	)

	sourceinterface := uint8(1)
	si := ie.NewInformationElement(
		ie.IESourceInterface,
		0,
		dt.OctetString([]byte{sourceinterface}),
	)

	//ueIPAddress, destinationIP address,
	ip := net.ParseIP("192.0.2.1")
	ueIPAddress := ie.NewUEIPAddress(false, true, true, false, ip, nil, 0)
	bb, err := ueIPAddress.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	ueIPAddressIE := ie.NewInformationElement(
		ie.IEUEIPaddress,
		0,
		dt.OctetString(bb),
	)

	pdi := ie.NewPDI(&si, nil, nil, &ueIPAddressIE, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	pdiB, err := pdi.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing PDI struct %+v", err)

	}

	pdiIE := ie.NewInformationElement(
		ie.IEPDI,
		0,
		dt.OctetString(pdiB),
	)

	createPDR := ie.NewCreatePDR(&pdrid, &precedence, &pdiIE, nil, &farID, nil, nil, nil)
	b, err := createPDR.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing CreatePDR %+v", err)

	}

	createPDRIE := ie.NewInformationElement(
		ie.IECreatePDR, //IEcode
		0,              //EntrepriseID
		dt.OctetString(b),
	)
	length := ie.IEBasicHeaderSize + createPDRIE.Len()

	//CreateFAR

	applyAction := dt.OctetString([]byte{0x02})
	a := ie.NewInformationElement(
		ie.IEApplyAction,
		0,
		applyAction,
	)

	destinationinterface := uint8(ie.Access)
	d := ie.NewInformationElement(
		ie.IEDestinationInterface,
		0,
		dt.OctetString([]byte{destinationinterface}),
	)

	ohcd := uint8(1)
	teid := uint32(100)
	ip4address := net.ParseIP("8.8.8.8")
	ohc := ie.NewOuterHeaderCreation(ohcd, teid, ip4address, nil, 0)
	b, err = ohc.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	ohcIE := ie.NewInformationElement(
		ie.IEOuterHeaderCreation,
		0,
		dt.OctetString(b),
	)

	fp := ie.NewForwardingParameters(&d, nil, nil, &ohcIE, nil, nil, nil, nil, nil)
	bb, err = fp.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	fpIE := ie.NewInformationElement(
		ie.IEForwardingParameters,
		0,
		dt.OctetString(bb),
	)

	createFAR := ie.NewCreateFAR(&farID, &a, &fpIE, nil, nil)
	bCreateFAR, err := createFAR.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)

	}

	createFARIE := ie.NewInformationElement(
		ie.IECreateFAR, //IEcode
		0,              //EntrepriseID
		dt.OctetString(bCreateFAR),
	)
	length = length + ie.IEBasicHeaderSize + createFARIE.Len()
	pfcpHeader := NewPFCPHeader(1, false, true, SessionModificationRequestType, length+12, sei, sn, 0)

	smr := NewPFCPSessionModificationRequest(pfcpHeader, nil, nil, nil, nil, nil, nil, nil, &createPDRIE, &createFARIE, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	bb, err = smr.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	ba := []byte{0x21, 0x34, 0x00, 0x5f,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x90, 0x00, 0x00, 0xc8, 0x00,
		0x00, 0x01, 0x00, 0x28, //CreatePDR
		0x00, 0x38, 0x00, 0x02, 0x00, 0x10, //IEPDRID
		0x00, 0x1d, 0x00, 0x04, 0x00, 0x00, 0x00, 0x0a, //IEPrecedence
		0x00, 0x02, 0x00, 0x0e, 0x00, 0x14, 0x00, 0x01, 0x01, 0x00, 0x5d, 0x00, 0x05, 0x06, 0xc0, 0x00, 0x02, 0x01, //IEPDI
		0x00, 0x6c, 0x00, 0x04, 0x00, 0x00, 0x00, 0x64, //IEFARID
		0x00, 0x03, 0x00, 0x23, //IECreateFAR
		0x00, 0x6c, 0x00, 0x04, 0x00, 0x00, 0x00, 0x64, //IEFARID
		0x00, 0x2c, 0x00, 0x01, 0x02, //IEApplyAction
		0x00, 0x04, 0x00, 0x12, //IEForwarding Parameters
		0x00, 0x2a, 0x00, 0x01, 0x00, //IEDestination Interface
		0x00, 0x54, 0x00, 0x09, 0x01, 0x00, 0x00, 0x00, 0x64, 0x08, 0x08, 0x08, 0X08, //IEOuterHeaderCreation
	}

	if !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}
}
