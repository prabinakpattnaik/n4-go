package qer

import (
	"bytes"
	"testing"
)

func TestBR(t *testing.T) {
	var u, d uint64
	u = 1024
	d = 512
	br := NewBR(u, d)
	b, err := br.Serialize()
	if err != nil {
		t.Fatalf("Error in BR serialization")
	}
	ba := []byte{0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00}

	if !bytes.Equal(b, ba) {
		t.Fatalf("unexpected value. want [%x], have [%x]", ba, b)
	}
}
func TestGateStatus(t *testing.T) {
	gs := NewGateStatus(OPEN, CLOSED)
	b, err := gs.Serialize()
	if err != nil {
		t.Fatalf("Error in GateStatus serialization")
	}
	ba := 0x01
	if b != byte(ba) {
		t.Fatalf("unexpected value. want [%x], have [%x]", ba, b)
	}
}

func TestIECreateQER(t *testing.T) {
	gs := NewGateStatus(OPEN, CLOSED)
	var u, d uint64
	u = 1024
	d = 512
	br := NewBR(u, d)
	createQER, err := NewCreateQER(1, 0, gs, br, br, 9, true)
	if err != nil {
		t.Fatalf("Error in Creating CreateQER")
	}
	b, err := createQER.Serialize()
	if err != nil {
		t.Fatalf("Error in Create QER serialization")
	}
	ba := []byte{0x00, 0x6d, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x19, 0x00, 0x01,
	             0x01, 0x00, 0x1a, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00,
		     0x02, 0x00, 0x00, 0x1b, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00,
		     0x00, 0x02, 0x00, 0x00, 0x7c, 0x00, 0x01, 0x09, 0x00, 0x7b, 0x00, 0x01, 0x09}
	if !bytes.Equal(b, ba) {
		t.Fatalf("unexpected value. want [%x], have [%x]", ba, b)
	}

}
