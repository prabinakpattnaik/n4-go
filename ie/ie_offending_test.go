package ie

import (
	"bytes"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	"testing"
)

func TestIEOffendingIE(t *testing.T) {
	b := []byte{0x10}

	i := NewInformationElement(
		IEOffendingIE,
		0,
		dt.OctetString(b),
	)

	if i.Length != uint16(len(b)) {
		t.Fatalf("Unexpected length. want %d, have %d", len(b), i.Length)
	}

	ba := []byte{0x00, 0x28, 0x00, 0x01, 0x10}

	b1, err := i.Serialize()

	if err != nil {
		t.Fatalf("Error in serializing %+v", err)

	}
	if !bytes.Equal(b1, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, b1)
	}

}
