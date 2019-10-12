package ie

import (
        "bytes"
        dt "github.com/fiorix/go-diameter/diam/datatype"
        "testing"
)


func TestNewPrecedence(t *testing.T) {
        pre := uint32(10)

        i := NewInformationElement(
                IEPrecedence,
                0,
                dt.Unsigned32(pre),
        )

        if i.Length != 4 {
                t.Fatalf("Unexpected length. want 4, have %d", i.Length)
        }

        bb, err := i.Serialize()
        if err != nil {
                t.Fatalf("Error in serializing %+v", err)
        }

        ba := []byte{0x00, 0x1D, 0x00, 0x04, 0x00, 0x00, 0x00, 0xa}

        if !bytes.Equal(bb, ba) {
                t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
        }

        t.Log(i)
}
