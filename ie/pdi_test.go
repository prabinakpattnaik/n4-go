package ie

import (
	"bytes"
	"net"
	"testing"

	"github.com/prabinakpattnaik/n4-go/util/util_3gpp"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)


func TestNewPDI(t *testing.T) {
	//included: source interface, Local F-TEID, Network Instance
	// not included: UE IP address, Traffic Endpoint ID, SDF Filter, Application ID
	//EthernetPDU Session Information, Ethernet Packet Filter, QFI, Framed-Route, Framed-Routing
	//Framed-IPv6-Route

	sourceinterface := uint8(1)

	si := NewInformationElement(
		IESourceInterface,
		0,
		dt.OctetString([]byte{sourceinterface}),
	)

	b, err := si.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	ft := []byte{0x01, 0x00, 0x00, 0xFF, 0xFF, 0xC0, 0xA8, 0x1, 0x65}

	fteid := NewInformationElement(
		IEFTEID,
		0,
		dt.OctetString(ft),
	)

	b1, err := fteid.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	ni := dt.OctetString("internet.mnc012.mcc345.gprs")

	networkInstance := NewInformationElement(
		IENetworkInstance,
		0,
		ni,
	)

	b2, err := networkInstance.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	b = append(b, b1...)
	b = append(b, b2...)

	i := NewInformationElement(
		IEPDI,
		0,
		dt.OctetString(b),
	)

	if i.Length != uint16(len(b)) {
		t.Fatalf("Unexpected length of FARID. want %d, have %d", (len(b)), i.Length)
	}

	ba := []byte{0x00, 0x02, 0x00, 0x31,
		0x00, 0x14, 0x00, 0x01, 0x01,
		0x00, 0x15, 0x00, 0x09, 0x01, 0x00, 0x00, 0xFF, 0xFF, 0xC0, 0xA8, 0x1, 0x65,
		0x00, 0x16, 0x00, 0x1b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x6d, 0x6e, 0x63, 0x30, 0x31, 0x32, 0x2e, 0x6d, 0x63, 0x63, 0x33, 0x34, 0x35, 0x2e, 0x67, 0x70, 0x72, 0x73}

	bb, err := i.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)

	}

	if !bytes.Equal(bb, ba) {

		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

	t.Log(i)
}



func TestNewPDIStruct(t *testing.T) {
	sourceinterface := uint8(1)

	si := NewInformationElement(
		IESourceInterface,
		0,
		dt.OctetString([]byte{sourceinterface}),
	)

	ft := []byte{0x01, 0x00, 0x00, 0xFF, 0xFF, 0xC0, 0xA8, 0x1, 0x65}
	fteid := NewInformationElement(
		IEFTEID,
		0,
		dt.OctetString(ft),
	)

	ni := dt.OctetString("internet.mnc012.mcc345.gprs")

	networkInstance := NewInformationElement(
		IENetworkInstance,
		0,
		ni,
	)
	ueIPaddress := NewUEIPAddress(false, true, false, false, net.ParseIP("8.8.8.8"), nil, 0)
	u, err := ueIPaddress.Serialize()

	ueIPAddressIE := NewInformationElement(
		IEUEIPaddress,
		0,
		dt.OctetString(u),
	)

	ba := []byte{0x00, 0x14, 0x00, 0x01, 0x01,
		0x00, 0x15, 0x00, 0x09, 0x01, 0x00, 0x00, 0xFF, 0xFF, 0xC0, 0xA8, 0x1, 0x65,
		0x00, 0x16, 0x00, 0x1b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2e, 0x6d, 0x6e, 0x63, 0x30, 0x31, 0x32, 0x2e, 0x6d, 0x63, 0x63, 0x33, 0x34, 0x35, 0x2e, 0x67, 0x70, 0x72, 0x73,
		0x00, 0x5d, 0x00, 0x05, 0x02, 0x08, 0x08, 0x08, 0x08}

	pdi := NewPDI(&si, &fteid, &networkInstance, &ueIPAddressIE, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	bb, err := pdi.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)

	}

	if !bytes.Equal(bb, ba) {

		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}
}


func TestNewPDIStructv1(t *testing.T) {
	si := &SourceInterface{
		InterfaceValue: SourceInterfaceAccess,
	}

	fteid := NewFTEID(true, false, false, false, 255, net.IPv4(192, 168, 1, 101), nil, 0)
	dnn := util_3gpp.Dnn("internet")
	ip := net.ParseIP("192.0.2.1")
	ueIPAddress := NewUEIPAddress(false, true, true, false, ip, nil, 0)

	pdi := PDI{
		SourceInterface: si,
		LocalFTEID:      fteid,
		NetworkInstance: &dnn,
		UEIPAddress:     ueIPAddress,
	}
	p, err := pdi.Serialize()
	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}
	/*ba := []byte{0x00, 0x02, 0x00, 0x28,
		0x00, 0x14, 0x00, 0x01, 0x00,
		0x00, 0x15, 0x00, 0x09, 0x01, 0x00, 0x00, 0x00, 0xFF, 0xC0, 0xA8, 0x1, 0x65,
		0x00, 0x16, 0x00, 0x09, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74,
		0x00, 0x5d, 0x00, 0x05, 0x06, 0xC0, 0x00, 0x02, 0x01}*/
	ba := []byte{0x00, 0x14, 0x00, 0x01, 0x00,
                0x00, 0x15, 0x00, 0x09, 0x01, 0x00, 0x00, 0x00, 0xFF, 0xC0, 0xA8, 0x1, 0x65,
                0x00, 0x16, 0x00, 0x09, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74,
                0x00, 0x5d, 0x00, 0x05, 0x06, 0xC0, 0x00, 0x02, 0x01}
	if !bytes.Equal(ba, p) {
		t.Fatalf("unexpected value of PDI IE. want [%x}, have [%x]", ba, p)
	}
}
