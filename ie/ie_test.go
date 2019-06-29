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
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

	t.Log(i)

}

func TestIENodeID(t *testing.T) {

	nodeID := []byte{0x0, 0xC0, 0xa8, 0x1, 0x65}
	i := NewInformationElement(
		IENodeID, //IEcode
		0,        //EntrepriseID
		dt.OctetString(nodeID),
	)

	if i.Length != 9 {
		t.Fatalf("Unexpected length. want 9, have %d", i.Length)
	}

	ba := []byte{0x00, 0x3c, 0x00, 0x9, 0x0, 0xC0, 0xa8, 0x1, 0x65}

	if b, _ := i.Serialize(); !bytes.Equal(b, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, b)
	}
}

func TestIEUPFunctionFeatures(t *testing.T) {

	i := NewInformationElement(
		IEUPFunctionFeatures,               //IEcode
		0,                                  //EntrepriseID
		dt.OctetString([]byte{0x00, 0X00}), //TODO Decode Up Function Features
	)

	if i.Length != 6 {
		t.Fatalf("Unexpected length. want 6, have %d", i.Length)
	}

	ba := []byte{0x00, 0x2B, 0x00, 0x06, 0x00, 0x00}

	if b, _ := i.Serialize(); !bytes.Equal(b, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, b)
	}

}

func TestIECPFunctionFeatures(t *testing.T) {

	i := NewInformationElement(
		IECPFunctionFeatures,         //IEcode
		0,                            //EntrepriseID
		dt.OctetString([]byte{0x00}), //TODO Decode Cp Function Features
	)

	if i.Length != 5 {
		t.Fatalf("Unexpected length. want 5, have %d", i.Length)
	}

	ba := []byte{0x00, 0x59, 0x00, 0x05, 0x00}

	if b, _ := i.Serialize(); !bytes.Equal(b, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, b)
	}

}
