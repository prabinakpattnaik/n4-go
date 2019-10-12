package ie

import (
        "bytes"
        dt "github.com/fiorix/go-diameter/diam/datatype"
        "testing"
)




func TestNewFARID(t *testing.T) {
        //Rule is dynamically provisioned by the CP Function
        farID := dt.Unsigned32(100)

        i := NewInformationElement(
                IEFARID,
                0,
                farID,
        )

        if i.Length != uint16(farID.Len()) {
                t.Fatalf("Unexpected length of FARID. want %d, have %d", farID.Len(), i.Length)
        }

        bb, err := i.Serialize()
        if err != nil {
                t.Fatalf("Error in serializing %+v", err)
        }

        ba := []byte{0x00, 0x6C, 0x00, 0x4, 0x00, 0x00, 0x00, 0x64}

        if !bytes.Equal(bb, ba) {
                t.Fatalf("unexpected value of FARID. want [%x}, have [%x]", ba, bb)
        }

        t.Log(i)
}

