package ie

import (
	"bytes"
	"net"
	"testing"
)

func TestNewFSEID(t *testing.T) {
	seid := uint64(500)
	ip := net.ParseIP("192.0.2.1")

	fseid := NewFSEID(true, false, seid, ip, nil)

	bb, err := fseid.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	ba := []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xF4, 0xC0, 0x00, 0x02, 0x01}

	if !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value of FSEID. want [%x}, have [%x]", ba, bb)
	}

	t.Log(fseid)
}
