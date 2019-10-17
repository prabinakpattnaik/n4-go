package msg

import (
	"bytes"
	"testing"
)

func TestNewPFCPHeader(t *testing.T) {
	sei := uint64(333)
	sn := uint32(444)
	version := uint8(1)
	mp := true
	s := true
	messagepriority := uint8(8)

	pfcpHeader := NewPFCPHeader(version, mp, s, HeartbeatRequestType, 12, sei, sn, messagepriority)
	ba := []byte{0x23, 0x1, 0x00, 0xc, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x4D, 0x00, 0x01, 0xBC, 0x80}

	if pfcpHeader.MessageType != HeartbeatRequestType {
		t.Fatalf("Unexpected message Type. want %d, hava %d", HeartbeatRequestType, pfcpHeader.MessageType)
	}

	if pfcpHeader.MessageLength != 0xc {
		t.Fatalf("Unexpected length. want 0xc, have %d", pfcpHeader.MessageLength)
	}

	if pfcpHeader.SessionEndpointIdentifier != sei {
		t.Fatalf("Unexpected session endpoint identidfier want %d, hava %d ", sei, pfcpHeader.SessionEndpointIdentifier)
	}

	if pfcpHeader.SequenceNumber != sn {
		t.Fatalf("Unexpected sequence number want %d, hava %d", sn, pfcpHeader.SequenceNumber)
	}

	if bb := pfcpHeader.Serialize(); !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

}

func TestNewPFCPHeaderWOS(t *testing.T) {
	sn := uint32(555)
	version := uint8(1)
	mp := false
	s := false

	pfcpHeader := NewPFCPHeader(version, mp, s, HeartbeatRequestType, 4, 0, sn, 0)
	ba := []byte{0x20, 0x1, 0x00, 0x4, 0x00, 0x02, 0x2B, 0x00}

	if pfcpHeader.MessageType != HeartbeatRequestType {
		t.Fatalf("Unexpected message Type. want %d, hava %d", HeartbeatRequestType, pfcpHeader.MessageType)
	}

	if pfcpHeader.MessageLength != 0x4 {
		t.Fatalf("Unexpected length. want 0xc, have %d", pfcpHeader.MessageLength)
	}

	if pfcpHeader.SessionEndpointIdentifier != 0 {
		t.Fatalf("Unexpected session endpoint identidfier want 0, hava %d ", pfcpHeader.SessionEndpointIdentifier)
	}

	if pfcpHeader.SequenceNumber != sn {
		t.Fatalf("Unexpected sequence number want %d, hava %d", sn, pfcpHeader.SequenceNumber)
	}

	if bb := pfcpHeader.Serialize(); !bytes.Equal(bb, ba) {
		t.Fatalf("unexpected value. want [%x}, have [%x]", ba, bb)
	}

}
