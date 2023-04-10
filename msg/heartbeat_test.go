package msg

import (
	"github.com/prabinakpattnaik/n4-go/ie"
	"bytes"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	"testing"
	"time"
)

func TestHeartbeatRequest(t *testing.T) {

	n := dt.Time(time.Unix(1377093974, 0))
	ba := []byte{0x20, 0x1, 0x00, 0x0C, 0x00, 0x00, 0x64, 0x00, 0x00, 0x60, 0x00, 0x04, 0xd5, 0xbf, 0x47, 0xd6}

	sn := uint32(100)
	pfcpHeader := NewPFCPHeader(1, false, false, HeartbeatRequestType, 12, 0, sn, 0)
	i := ie.NewInformationElement(
		ie.IERecoveryTimestamp,
		0,
		n,
	)

	heartbeat := NewHeartbeat(pfcpHeader, &i)
	if bb, _ := heartbeat.Serialize(); !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

}

func TestHeartbeatResponse(t *testing.T) {

	n := dt.Time(time.Unix(1377093974, 0))
	ba := []byte{0x20, 0x2, 0x00, 0x0C, 0x00, 0x00, 0x64, 0x00, 0x00, 0x60, 0x00, 0x04, 0xd5, 0xbf, 0x47, 0xd6}

	sn := uint32(100)
	pfcpHeader := NewPFCPHeader(1, false, false, HeartbeatResponseType, 12, 0, sn, 0)
	i := ie.NewInformationElement(
		ie.IERecoveryTimestamp,
		0,
		n,
	)

	heartbeat := NewHeartbeat(pfcpHeader, &i)
	if bb, _ := heartbeat.Serialize(); !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

}
