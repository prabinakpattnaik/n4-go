package ie

import (
	"bytes"
	"net"
	"testing"

	dt "github.com/fiorix/go-diameter/diam/datatype"
)

func TestIEOuterHeaderCreation(t *testing.T) {
	// GTP-U/UDP/IPv4 (5/1)
	//TEID
	outerHeaderCreation := dt.OctetString([]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0xff, 0x08, 0x08, 0x08, 0x8})

	i := NewInformationElement(
		IEOuterHeaderCreation,
		0,
		outerHeaderCreation,
	)

	if i.Length != uint16(outerHeaderCreation.Len()) {
		t.Fatalf("OuterHeaderCreation:Unexpected length. want %d, have %d", outerHeaderCreation.Len(), i.Length)
	}

	bb, err := i.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	ba := []byte{0x00, 0x54, 0x00, 0x0a, 0x01, 0x00, 0x00, 0x00, 0x00, 0xff, 0x08, 0x08, 0x08, 0x08}

	if !bytes.Equal(bb, ba) {
		t.Fatalf("outerHeaderCreation:unexpected value. want [%x}, have [%x]", ba, bb)
	}

	t.Log(i)

}

func TestNewOuterHeaderCreation(t *testing.T) {
	ohcd := uint8(1)
	teid := uint32(100)
	ip4address := net.ParseIP("8.8.8.8")

	ohc := NewOuterHeaderCreation(ohcd, teid, ip4address, nil, 0)
	b, err := ohc.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}
	a := []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x64, 0x08, 0x08, 0x08, 0x08}

	if !bytes.Equal(b, a) {
		t.Fatalf("outerHeaderCreation:unexpected value. want [%x}, have [%x]", a, b)
	}

}
