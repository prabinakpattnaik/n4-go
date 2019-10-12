package ie

import (
        "bytes"
        dt "github.com/fiorix/go-diameter/diam/datatype"
        "testing"
)



func TestNewOuterHeaderRemoval(t *testing.T) {
        //GTP-U/UDP/IPv4 store GTP-U extension header, store GTP-U message type

        outerHeaderRemoval := dt.OctetString([]byte{0x00})

        i := NewInformationElement(
                IEOuterHeaderRemoval,
                0,
                outerHeaderRemoval,
        )

        if i.Length != uint16(outerHeaderRemoval.Len()) {
                t.Fatalf("Unexpected length. want %d, have %d", outerHeaderRemoval.Len(), i.Length)
        }

        bb, err := i.Serialize()
        if err != nil {
                t.Fatalf("Error in serializing %+v", err)
        }

        ba := []byte{0x00, 0x5F, 0x00, 0x1, 0x00}

        if !bytes.Equal(bb, ba) {
                t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
        }

        t.Log(i)
}

