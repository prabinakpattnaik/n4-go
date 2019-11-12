package ie

import (
	"bytes"
	"testing"

	dt "github.com/fiorix/go-diameter/diam/datatype"
)

func TestCreatedPDR(t *testing.T) {
	ruleID := []byte{0x00, 0x10}
	pdrid := NewInformationElement(
		IEPDRID,
		0,
		dt.OctetString(ruleID),
	)
	createdPDR := NewCreatedPDR(&pdrid, nil)
	b, err := createdPDR.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing CreatedPDR %+v", err)

	}
	ba := []byte{0x00, 0x38, 0x00, 0x02, 0x00, 0x10}
	if !bytes.Equal(b, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, b)

	}

}
