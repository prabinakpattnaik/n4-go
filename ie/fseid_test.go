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

func TestNewFSEIDFromBytes(t *testing.T) {
	b := []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xac, 0x11, 0x00, 0x05}
	fseid := NewFSEIDFromByte(b)
	ipaddress := []byte{0xac, 0x11, 0x00, 0x05}
	if fseid.SEID != 2 {
		t.Fatalf("wrong SEID %d", fseid.SEID)
	}
	if !bytes.Equal(fseid.IP4Address, ipaddress) {
		t.Fatalf("unexpected value of FSEID IPV4 address. want [%x}, have [%x]", ipaddress, fseid.IP4Address)
	}

}
