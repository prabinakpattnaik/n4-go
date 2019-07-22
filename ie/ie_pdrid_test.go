package ie

import (
	"bytes"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	"testing"
)

func TestIEPDRID(t *testing.T) {

	i := NewInformationElement(

		IEPDRID, //IEcode

		0, //EntrepriseID

		dt.OctetString([]byte{0x00, 0x02}), //RuleID  0x0002

	)

	if i.Length != 2 {

		t.Fatalf("Unexpected length. want 2, have %d", i.Length)

	}

	ba := []byte{0x00, 0x38, 0x00, 0x02, 0x00, 0x02}

	if b, _ := i.Serialize(); !bytes.Equal(b, ba) {

		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, b)

	}

}
