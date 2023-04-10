package ie

import (
	"bytes"
	"net"
	"testing"
)

func TestUPIPResourceFunction(t *testing.T) {
	data := []byte{0x29, 0x80, 0xac, 0x13, 0x00, 0x02, 0x02, 0x63, 0x70}
	ipaddress := []byte{0xac, 0x13, 0x00, 0x02}

	upIPResourceInformation := NewUPIPResourceInformationFromByte(9, data)
	if !bytes.Equal(upIPResourceInformation.IPv4Address, ipaddress) {
		t.Fatalf("unexpected value. want [%x], have [%x]", ipaddress, upIPResourceInformation.IPv4Address)

	}
	if !upIPResourceInformation.ASSONI {
		t.Fatalf("unexpected value. want true, have %t", upIPResourceInformation.ASSONI)
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
	/*
	upIPResourceInformation = NewUPIPResourceInformationFromByte(uint16(len(bb)), bb)
	if !bytes.Equal(upIPResourceInformation.IPv4Address, ipv4address.To4()) {
		t.Fatalf("unexpected value. want [%v], have [%v]", ipv4address, upIPResourceInformation.IPv4Address)
	}*/

}
