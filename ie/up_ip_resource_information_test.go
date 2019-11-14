package ie

import (
	"bytes"
	"net"
	"testing"
)

func TestUPIPResourceFunction(t *testing.T) {
	data := []byte{0x11, 0x0f, 0xC0, 0xa8, 0x1, 0x65}
	ipaddress := []byte{0xC0, 0xa8, 0x1, 0x65}

	upIPResourceInformation := NewUPIPResourceInformationFromByte(6, data)
	if !bytes.Equal(upIPResourceInformation.IPv4Address, ipaddress) {
		t.Fatalf("unexpected value. want [%x], have [%x]", ipaddress, upIPResourceInformation.IPv4Address)

	}

}

func TestUPIPResourceFunctionStruct(t *testing.T) {
	ipv4address := net.IPv4(8, 8, 8, 8)
	upIPResourceInformation := NewUPIPResourceInformation(true, false, 0, false, false, 0, ipv4address, nil, nil, 0)
	ba, err := upIPResourceInformation.Serialize()
	if err != nil {
		t.Fatalf("Error is popped up %+v\n", err)
	}
	bb := []byte{0x01, 0x00, 0x08, 0x08, 0x08, 0x08}
	if !bytes.Equal(ba, bb) {
		t.Fatalf("unexpected value. want [%x], have [%x]", bb, ba)

	}
	upIPResourceInformation = NewUPIPResourceInformationFromByte(uint16(len(ba)), ba)
	if !bytes.Equal(upIPResourceInformation.IPv4Address, ipv4address.To4()) {
		t.Fatalf("unexpected value. want [%v], have [%v]", ipv4address, upIPResourceInformation.IPv4Address)
	}

}
