package session

import (
	"bytes"
	"net"
	"testing"

	"bitbucket.org/sothy5/n4-go/ie"
)

func TestSessionManagement(t *testing.T) {
	sei := uint64(400)
	sn := uint32(200)
	seid := uint64(1000)
	nodeIP := net.ParseIP("192.168.1.101")
	pdrid := uint16(10)
	farid := uint32(100)
	sourceinterface := uint8(1)
	fteid := ie.NewFTEID(true, false, false, false, 255, nodeIP, nil, 0)
	aa := uint8(2)
	destionationinterface := uint8(1)
	pfcpSessionEstablishmentRequest, err := CreateSession(sei, sn, nodeIP, seid, pdrid, farid, sourceinterface, fteid, aa, destionationinterface)
	if err != nil {
		t.Fatalf("Create New Session Error %+v\n", err)
	}
	b, err := pfcpSessionEstablishmentRequest.Serialize()
	if err != nil {
		t.Fatalf("New PFCP Session serialization Error %+v\n", err)
	}

	bb := []byte{0x21, 0x32, 0x00, 0x6c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x90, 0x00, 0x00, 0xc8, 0x00,
		0x00, 0x3c, 0x00, 0x05, 0x00, 0xc0, 0xa8, 0x01, 0x65, //NODEID
		0x00, 0x39, 0x00, 0x0d, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xe8, 0xc0, 0xa8, 0x01, 0x65, //FTEID
		0x00, 0x01, 0x00, 0x28, 0x00, 0x38, 0x00, 0x02, 0x00, 0x0a, //CreatePDR, PDRID
		0x00, 0x1d, 0x00, 0x04, 0x00, 0x00, 0x00, 0x0a, //Precedance
		0x00, 0x02, 0x00, 0x09, 0x00, 0x14, 0x00, 0x01, 0x01, 0x00, 0x15, 0x00, 0x00, //PDI
		0x00, 0x5f, 0x00, 0x01, 0x06, //OuterHeaderRemoval
		0x00, 0x6c, 0x00, 0x04, 0x00, 0x00, 0x00, 0x64, //FARID
		0x00, 0x03, 0x00, 0x16, 0x00, 0x6c, 0x00, 0x04, 0x00, 0x00, 0x00, 0x64, //CreateFAR,FARID,
		0x00, 0x2c, 0x00, 0x01, 0x02, //ApplyAction
		0x00, 0x04, 0x00, 0x05, 0x00, 0x2a, 0x00, 0x01, 0x01} //Forwarding Parameters

	if !bytes.Equal(bb, b) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", bb, b)
	}

}
