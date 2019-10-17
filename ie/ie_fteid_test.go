package ie

import (
	"bytes"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	"testing"
)

func TestIEFTEID(t *testing.T) {
	//CP assigns IPv4 address and TEID (Ox00,0x00,OxFF,0xFF)
	// IPv4 (192.168.1.101)
	b := []byte{0x01, 0x00, 0x00, 0xFF, 0xFF, 0xC0, 0xA8, 0x1, 0x65}

	i := NewInformationElement(
		IEFTEID,
		0,
		dt.OctetString(b),
	)

	if i.Length != uint16(len(b)) {
		t.Fatalf("Unexpected length. want %d, have %d", len(b), i.Length)
	}

	ba := []byte{0x00, 0x15, 0x00, 0x09, 0x01, 0x00, 0x00, 0xFF, 0xFF, 0xC0, 0xA8, 0x1, 0x65}

	b1, err := i.Serialize()

	if err != nil {
		t.Fatalf("Error in serializing %+v", err)

	}
	if !bytes.Equal(b1, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, b1)
	}

}
