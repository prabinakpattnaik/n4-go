package ie

import (
	"bytes"
	"testing"
	"time"

	dt "github.com/fiorix/go-diameter/diam/datatype"
)

var testIE = [][]byte{
	//recovery-timestamp
	{},
}

func TestNewInformationElement(t *testing.T) {

	i := NewInformationElement(
		IERecoveryTimestamp, //IEcode
		0,                   //EntrepriseID
		dt.Time(time.Now()), //seconds in since January 1, 1900 UTC
	)

	if i.Length != 8 {
		t.Fatalf("Unexpected length. want 8, have %d", i.Length)
	}

	_, err := i.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	n := dt.Time(time.Unix(1377093974, 0))
	ba := []byte{0x00, 0x60, 0x00, 0x08, 0xd5, 0xbf, 0x47, 0xd6}

	ii := NewInformationElement(
		IERecoveryTimestamp,
		0,
		n,
	)
	if ii.Length != 8 {
		t.Fatalf("Unexpected length in ii, want 8, have %d", ii.Length)
	}
	if bb, _ := ii.Serialize(); !bytes.Equal(bb, ba) {
		//t.Fatalf("[%x]", n.Serialize())
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

	t.Log(i)

}
