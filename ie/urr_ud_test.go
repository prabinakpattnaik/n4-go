package ie

import (
	"bytes"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	"testing"
)

func TestNewURRID(t *testing.T) {
	//Rule is predefined in the UP Function.
	urrID := dt.Unsigned32(2147483663)

	i := NewInformationElement(
		IEURRID,
		0,
		urrID,
	)

	if i.Length != uint16(urrID.Len()) {
		t.Fatalf("Unexpected length of URRID. want %d, have %d", urrID.Len(), i.Length)
	}

	bb, err := i.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	ba := []byte{0x00, 0x51, 0x00, 0x4, 0x80, 0x00, 0x00, 0x0f}

	if !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value of URRID. want [%x}, have [%x]", ba, bb)
	}

	t.Log(i)
}
