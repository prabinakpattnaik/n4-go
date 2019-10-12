package ie

import (
        "bytes"
        dt "github.com/fiorix/go-diameter/diam/datatype"
        "testing"
)


func TestNewSourceInterface(t *testing.T) {
        sourceinterface := 1
        b := uint8(0xFF & sourceinterface)

        i := NewInformationElement(
                IESourceInterface,
                0,
                dt.OctetString([]byte{b}),
        )

        if i.Length != 1 {
                t.Fatalf("Unexpected length. want 1, have %d", i.Length)
        }

        bb, err := i.Serialize()
        if err != nil {
                t.Fatalf("Error in serializing %+v", err)
        }

        ba := []byte{0x00, 0x14, 0x00, 0x01, 0x01}

        if !bytes.Equal(bb, ba) {
                t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
        }

        t.Log(i)
}

