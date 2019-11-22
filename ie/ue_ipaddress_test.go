package ie

import (
	"bytes"
	"net"
	"testing"
)

func TestUEIPAddressV4(t *testing.T) {
	ip := net.ParseIP("192.0.2.1")
	ueIPAddress := NewUEIPAddress(false, true, true, false, ip, nil, 0)
	bb, err := ueIPAddress.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	ba := []byte{0x06, 0xC0, 0x00, 0x02, 0x01}

	if !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value of FSEID. want [%x}, have [%x]", ba, bb)
	}
}
