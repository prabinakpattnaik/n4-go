package ie

import (
	"bytes"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	"testing"
)

func TestNewOuterHeaderCreation(t *testing.T) {
	// GTP-U/UDP/IPv4 (5/1)
	//TEID
	outerHeaderCreation := dt.OctetString([]byte{0x01, 0x00, 0x00, 0x00, 0xff})

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

	ba := []byte{0x00, 0x54, 0x00, 0x5, 0x01, 0x00, 0x00, 0x00, 0xff}

	if !bytes.Equal(bb, ba) {
		t.Fatalf("outerHeaderCreation:unexpected value. want [%x}, have [%x]", ba, bb)
	}

	t.Log(i)

}
