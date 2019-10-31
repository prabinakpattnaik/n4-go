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

func TestMessageFromBytes(t *testing.T) {
	ba := []byte{0x20, 0x05, 0x00, 0x25, 0x00, 0x00, 0xc8, 0X00, 0x00, 0x3c, 0x00, 0x5, 0x0, 0xC0, 0xa8, 0x1, 0x65, 0x00, 0x60, 0x00, 0x04, 0xd5, 0xbf, 0x47, 0xd6, 0x00, 0x2B, 0x00, 0x02, 0x00, 0x00, 0x00, 0x74, 0x00, 0x06, 0x11, 0x0F, 0xC0, 0xa8, 0x1, 0x65}
	pfcpMessage, err := MessageFromBytes(ba)

	if err != nil {
		t.Fatalf("Error in serializing %+v", err)
	}

	if pfcpMessage.Header.Version != 1 {
		t.Fatalf("unexpected version. want 1, have %d", pfcpMessage.Header.Version)
	}

	if pfcpMessage.Header.MP != false {
		t.Fatalf("unexpected MP value. want false, have %T", pfcpMessage.Header.MP)
	}

	if pfcpMessage.Header.S != false {
		t.Fatalf("unexpected S value. want false, have %T", pfcpMessage.Header.S)
	}

	if pfcpMessage.Header.MessageType != AssociationSetupRequestType {
		t.Fatalf("unexpected MessageType value. want %d, have %d", AssociationSetupRequestType, pfcpMessage.Header.MessageType)
	}

	if pfcpMessage.Header.MessageLength != 37 {
		t.Fatalf("unexpected version. want 37, have %d", pfcpMessage.Header.MessageLength)
	}

	if pfcpMessage.Header.SequenceNumber != 200 {
		t.Fatalf("unexpected SequenceNumber value. want 200, have %d", pfcpMessage.Header.SequenceNumber)
	}

	FromPFCPMessage(pfcpMessage)
	//TODO check the return variable of above function

}
