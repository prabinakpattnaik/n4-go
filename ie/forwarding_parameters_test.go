package ie

import (
	"bytes"
	"net"
	"testing"

	dt "github.com/fiorix/go-diameter/diam/datatype"
)

func TestIEForwardingParameters(t *testing.T) {
	//Included Destination interface
	//not included network instance, redirect information, outer header creation, transport level marking, forwarding policy, HeaderEnrichment
	//linked traffic Endpoint ID, proxying

	destinationinterface := uint8(1)

	i := NewInformationElement(
		IEDestinationInterface,
		0,
		dt.OctetString([]byte{destinationinterface}),
	)
	b, err := i.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	fp := NewInformationElement(

		IEForwardingParameters,
		0,
		dt.OctetString(b),
	)

	if fp.Length != uint16(len(b)) {

		t.Fatalf("Unexpected length. want %d, have %d", len(b), fp.Length)

	}

	ba := []byte{0x00, 0x04, 0x00, 0x5, 0x00, 0x2a, 0x00, 0x01, 0x01}

	bb, err := fp.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)

	}

	if !bytes.Equal(bb, ba) {

		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)

	}

}

func TestForwardingParameters(t *testing.T) {
	destinationinterface := uint8(1)
	d := NewInformationElement(
		IEDestinationInterface,
		0,
		dt.OctetString([]byte{destinationinterface}),
	)

	ohcd := uint8(1)
	teid := uint32(100)
	ip4address := net.ParseIP("8.8.8.8")
	ohc := NewOuterHeaderCreation(ohcd, teid, ip4address, nil, 0)
	b, err := ohc.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	ohcIE := NewInformationElement(
		IEOuterHeaderCreation,
		0,
		dt.OctetString(b),
	)

	fp := NewForwardingParameters(&d, nil, nil, &ohcIE, nil, nil, nil, nil, nil)
	ba := []byte{0x00, 0x2a, 0x00, 0x01, 0x01, 0x00, 0x54, 0x00, 0x09, 0x01, 0x00, 0x00, 0x00, 0x64, 0x08, 0x08, 0x08, 0x08}

	bb, err := fp.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	if !bytes.Equal(bb, ba) {

		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)

	}

}
