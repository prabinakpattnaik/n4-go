package ie

import (
        "bytes"
        dt "github.com/fiorix/go-diameter/diam/datatype"
        "testing"
)

func TestNewPDRIDIE(t *testing.T) {
        ruleID := []byte{0x00, 0x10}
        i := NewInformationElement(
                IEPDRID,
                0,
                dt.OctetString(ruleID),
        )

        if i.Length != 2 {
                t.Fatalf("Unexpected length. want 2, have %d", i.Length)
        }

        bb, err := i.Serialize()
        if err != nil {
                t.Fatalf("Error in serializing %+v", err)
        }

        ba := []byte{0x00, 0x38, 0x00, 0x02, 0x00, 0x10}

        if !bytes.Equal(bb, ba) {
                t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
        }

        t.Log(i)

}

