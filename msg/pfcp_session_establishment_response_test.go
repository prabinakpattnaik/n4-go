package msg

import (
	"bytes"
	"net"
	"testing"

	"bitbucket.org/sothy5/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

func TestPFCPSessionEstablishmentResponse(t *testing.T) {
	var length uint16
	nodeID := []byte{0x0, 0xC0, 0xa8, 0x1, 0x65}
	n := ie.NewInformationElement(
		ie.IENodeID,
		0,
		dt.OctetString(nodeID),
	)

	length = ie.IEBasicHeaderSize + n.Len()

	c := ie.NewInformationElement(
		ie.IECause,
		0,
		dt.OctetString([]byte{0x01}),
	)
	length = length + ie.IEBasicHeaderSize + c.Len()

	seid := uint64(1000)
	ip4address := net.ParseIP("192.168.1.101")
	fseid := ie.NewFSEID(true, false, seid, ip4address, nil)
	bb, err := fseid.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing of FSEID %+v", err)
	}
	upfseid := ie.NewInformationElement(
		ie.IEFSEID,
		0,
		dt.OctetString(bb),
	)

	length = length + ie.IEBasicHeaderSize + upfseid.Len()
	t.Logf("upfseid length %d\n", length)

	//CreatedPDR
	ruleID := []byte{0x00, 0x10}
	pdrid := ie.NewInformationElement(
		ie.IEPDRID,
		0,
		dt.OctetString(ruleID),
	)

	createdPDR := ie.NewCreatedPDR(&pdrid, nil)
	b, err := createdPDR.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing CreatedPDR %+v", err)

	}

	createdPDRIE := ie.NewInformationElement(
		ie.IECreatedPDR,
		0,
		dt.OctetString(b),
	)
	length = length + ie.IEBasicHeaderSize + createdPDRIE.Len()

	sei := uint64(400)
	sn := uint32(200)
	pfcpHeader := NewPFCPHeader(1, false, true, SessionEstablishmentResponseType, length+12, sei, sn, 0)

	t.Logf("length %d\n", length)

	sr := NewPFCPSessionEstablishmentResponse(pfcpHeader, &n, &c, nil, &upfseid, &createdPDRIE, nil, nil, nil, nil)
	ba := []byte{0x21, 0x33, 0x00, 0x35, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x90, 0x00, 0x00, 0xc8, 0x00,
		0x00, 0x3c, 0x00, 0x5, 0x0, 0xC0, 0xa8, 0x1, 0x65,
		0x00, 0x13, 0x00, 0x01, 0x01,
		0x00, 0x39, 0x00, 0x0d, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xe8, 0xc0, 0xa8, 0x01, 0x65,
		0x00, 0x08, 0x00, 0x06,
		0x00, 0x38, 0x00, 0x02, 0x00, 0x10,
	}

	if bb, _ := sr.Serialize(); !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

}
