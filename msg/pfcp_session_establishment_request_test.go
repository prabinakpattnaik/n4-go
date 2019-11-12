package msg

import (
	"bytes"
	"net"
	"testing"

	"bitbucket.org/sothy5/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

func TestPFCPSessionEstablishmentRequest(t *testing.T) {
	sei := uint64(400)
	sn := uint32(200)
	pfcpHeader := NewPFCPHeader(1, false, true, SessionEstablishmentRequestType, 127+12, sei, sn, 0)

	var length uint16
	nodeID := []byte{0x0, 0xC0, 0xa8, 0x1, 0x65}
	n := ie.NewInformationElement(
		ie.IENodeID, //IEcode
		0,           //EntrepriseID
		dt.OctetString(nodeID),
	)

	length = ie.IEBasicHeaderSize + n.Len()

	seid := uint64(1000)
	ip4address := net.ParseIP("192.168.1.101")
	fseid := ie.NewFSEID(true, false, seid, ip4address, nil)
	bb, err := fseid.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing of FSEID %+v", err)
	}
	cpfseid := ie.NewInformationElement(
		ie.IEFSEID, //IEcode
		0,          //EntrepriseID
		dt.OctetString(bb),
	)

	length = length + ie.IEBasicHeaderSize + cpfseid.Len()

	//CreatePDR

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

	oHR := dt.OctetString([]byte{0x00})
	outerHeaderRemoval := ie.NewInformationElement(
		ie.IEOuterHeaderRemoval,
		0,
		oHR,
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

	ft := []byte{0x01, 0x00, 0x00, 0xFF, 0xFF, 0xC0, 0xA8, 0x1, 0x65}
	fteid := ie.NewInformationElement(
		ie.IEFTEID,
		0,
		dt.OctetString(ft),
	)

	ni := dt.OctetString("internet.mnc012.mcc345.gprs")
	networkInstance := ie.NewInformationElement(
		ie.IENetworkInstance,
		0,
		ni,
	)
	pdi := ie.NewPDI(&si, &fteid, &networkInstance, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	pdiB, err := pdi.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing PDI struct %+v", err)

	}

	pdiIE := ie.NewInformationElement(
		ie.IEPDI,
		0,
		dt.OctetString(pdiB),
	)

	createPDR := ie.NewCreatePDR(&pdrid, &precedence, &pdiIE, &outerHeaderRemoval, &farID, nil, nil, nil)
	b, err := createPDR.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing CreatePDR %+v", err)

	}

	createPDRIE := ie.NewInformationElement(
		ie.IECreatePDR, //IEcode
		0,              //EntrepriseID
		dt.OctetString(b),
	)
	length = length + ie.IEBasicHeaderSize + createPDRIE.Len()

	//CreateFAR

	applyAction := dt.OctetString([]byte{0x02})
	a := ie.NewInformationElement(
		ie.IEApplyAction,
		0,
		applyAction,
	)

	createFAR := ie.NewCreateFAR(&farID, &a, nil, nil, nil)
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

	sr := NewPFCPSessionEstablishmentRequest(pfcpHeader, &n, &cpfseid, &createPDRIE, &createFARIE, nil, nil, nil, nil, nil, nil, nil, nil)

	ba := []byte{0x21, 0x32, 0x00, 0x8b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x90, 0x00, 0x00, 0xc8, 0x00,
		0x00, 0x3c, 0x00, 0x05, 0x00, 0xc0, 0xa8, 0x01, 0x65,
		0x00, 0x39, 0x00, 0x0d, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xe8, 0xc0, 0xa8, 0x01, 0x65,
		0x00, 0x01, 0x00, 0x50, 0x00, 0x38, 0x00, 0x02, 0x00, 0x10, 0x00, 0x1d, 0x00, 0x04,
		0x00, 0x00, 0x00, 0x0a, 0x00, 0x02, 0x00, 0x31, 0x00, 0x14, 0x00, 0x01, 0x01, 0x00,
		0x15, 0x00, 0x09, 0x01, 0x00, 0x00, 0xff, 0xff, 0xc0, 0xa8, 0x01, 0x65, 0x00, 0x16,
		0x00, 0x1b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x6d, 0x6e, 0x63,
		0x30, 0x31, 0x32, 0x2e, 0x6d, 0x63, 0x63, 0x33, 0x34, 0x35, 0x2e, 0x67, 0x70, 0x72,
		0x73, 0x00, 0x5f, 0x00, 0x01, 0x00, 0x00, 0x6c, 0x00, 0x04, 0x00, 0x00, 0x00, 0x64,

		0x00, 0x03, 0x00, 0x0d, 0x00, 0x6c, 0x00, 0x04, 0x00, 0x00, 0x00, 0x64, 0x00, 0x2c,
		0x00, 0x01, 0x02,
	}

	if bb, _ := sr.Serialize(); !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

}
