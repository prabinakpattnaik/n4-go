package ie

import (
	"bytes"
	"testing"
)

func TestUPIPResourceFunction(t *testing.T) {
	data := []byte{0x11, 0x0f, 0xC0, 0xa8, 0x1, 0x65}
	ipaddress := []byte{0xC0, 0xa8, 0x1, 0x65}

	upIPResourceInformation := NewUPIPResourceInformation(6, data)
	if !bytes.Equal(upIPResourceInformation.IPv4Address, ipaddress) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ipaddress, upIPResourceInformation.IPv4Address)

	}

}
