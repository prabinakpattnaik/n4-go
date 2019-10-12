package ie

import (
	"bytes"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	"testing"
)

func TestNewApplyActionIE(t *testing.T) {
	//Forward Action
	applyAction := dt.OctetString([]byte{0x02})
	i := NewInformationElement(
		IEApplyAction,
		0,
		applyAction,
	)

	if i.Length != uint16(applyAction.Len()) {
		t.Fatalf("Unexpected length. want %d, have %d", applyAction.Len(), i.Length)
	}

	bb, err := i.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	ba := []byte{0x00, 0x2C, 0x00, 0x01, 0x02}

	if !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

	t.Log(i)

}
