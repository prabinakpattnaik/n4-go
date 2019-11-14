package ie

import (
	"bytes"
	"net"
	"testing"
)

func TestFTEID(t *testing.T) {
	//CP assigns IPv4 address and TEID (Ox00,0x00,Ox00,0xFF)
	// IPv4 (192.168.1.101)

	fteid := NewFTEID(true, false, false, false, 255, net.IPv4(192, 168, 1, 101), nil, 0)

	ba := []byte{0x01, 0x00, 0x00, 0x00, 0xFF, 0xC0, 0xA8, 0x1, 0x65}

	bb, err := fteid.Serialize()

	if err != nil {
		t.Fatalf("Error in serializing %+v", err)

	}
	if !bytes.Equal(ba, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

}
