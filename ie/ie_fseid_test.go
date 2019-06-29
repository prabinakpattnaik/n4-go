package ie

import (
	"bytes"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	"testing"
)

func TestIEFSEID(t *testing.T) {

	i := NewInformationElement(
		IEFSEID, //IEcode
		0,       //EntrepriseID
		dt.OctetString([]byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xC0, 0xa8, 0x1, 0x65}), //TODO Decode FSEID (192.168.1.101), SEID 0x0000000000000001
	)

	if i.Length != 13 {
		t.Fatalf("Unexpected length. want 13, have %d", i.Length)
	}

	ba := []byte{0x00, 0x39, 0x00, 0x0d, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xC0, 0xa8, 0x1, 0x65}

	if b, _ := i.Serialize(); !bytes.Equal(b, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, b)
	}

}
