package ie

import (
	"bytes"
	"testing"

	dt "github.com/fiorix/go-diameter/diam/datatype"
)

func TestCause(t *testing.T) {
	i := NewInformationElement(
		IECause,                      //IEcode
		0,                            //EntrepriseID
		dt.OctetString([]byte{0x01}), //TODO Decode Request Accepted code
	)

	if i.Length != 1 {
		t.Fatalf("Unexpected length. want 6, have %d", i.Length)
	}

	ba := []byte{0x00, 0x13, 0x00, 0x01, 0x01}

	if b, _ := i.Serialize(); !bytes.Equal(b, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, b)
	}

}
