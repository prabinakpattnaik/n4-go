package msg

import (
	"bitbucket.org/sothy5/n4-go/ie"
	"bytes"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	"testing"
	"time"
)

func TestPFCPAssociationSetupRequest(t *testing.T) {
	sn := uint32(200)
	pfcpHeader := NewPFCPHeader(1, false, false, AssociationSetupRequestType, 37, 0, sn, 0)

	nodeID := []byte{0x0, 0xC0, 0xa8, 0x1, 0x65}
	n := ie.NewInformationElement(
		ie.IENodeID, //IEcode
		0,           //EntrepriseID
		dt.OctetString(nodeID),
	)

	td := dt.Time(time.Unix(1377093974, 0))

	r := ie.NewInformationElement(
		ie.IERecoveryTimestamp,
		0,
		td,
	)

	u := ie.NewInformationElement(
		ie.IEUPFunctionFeatures,
		0,
		dt.OctetString([]byte{0x00, 0X00}),
	)

	ui := ie.NewInformationElement(
		ie.IEUserPlaneIPResourceInformation,
		0,
		dt.OctetString([]byte{0x11, 0x0F, 0xC0, 0xa8, 0x1, 0x65}),
	)

	ar := NewPFCPAssociationSetupRequest(pfcpHeader, n, r, u, nil, ui)

	ba := []byte{0x20, 0x05, 0x00, 0x25, 0x00, 0x00, 0xc8, 0X00, 0x00, 0x3c, 0x00, 0x5, 0x0, 0xC0, 0xa8, 0x1, 0x65, 0x00, 0x60, 0x00, 0x04, 0xd5, 0xbf, 0x47, 0xd6, 0x00, 0x2B, 0x00, 0x02, 0x00, 0x00, 0x00, 0x74, 0x00, 0x06, 0x11, 0x0F, 0xC0, 0xa8, 0x1, 0x65}
	if bb := ar.Serialize(); !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

}