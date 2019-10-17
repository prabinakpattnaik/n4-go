package ie

import (
        "bytes"
        dt "github.com/fiorix/go-diameter/diam/datatype"
        "testing"
)



func TestNewNetworkInstance(t *testing.T) {
        networkinstance := dt.OctetString("internet.mnc012.mcc345.gprs")

        i := NewInformationElement(
                IENetworkInstance,
                0,
                networkinstance,
        )

        if i.Length != uint16(networkinstance.Len()) {
                t.Fatalf("Unexpected length. want %d, have %d", networkinstance.Len(), i.Length)
        }

        bb, err := i.Serialize()
        if err != nil {
                t.Fatalf("Error in serializing %+v", err)
        }

        ba := []byte{0x00, 0x16, 0x00, 0x1b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x6d, 0x6e, 0x63, 0x30, 0x31, 0x32, 0x2e, 0x6d, 0x63, 0x63, 0x33, 0x34, 0x35, 0x2e, 0x67, 0x70, 0x72, 0x73}

        if !bytes.Equal(bb, ba) {
                t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
        }

        t.Log(i)
}
