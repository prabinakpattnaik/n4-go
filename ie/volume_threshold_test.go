package ie

import (
	"testing"
)

func TestVolume(t *testing.T) {
	b := []byte{0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	v := Volume{}
	err := v.VolumeFromBytes(b)
	if err != nil {
		t.Fatalf("Volume is not valid %+v\n", err)
	}
	if !v.Tovol {
		t.Fatalf("Tovol flag is wrong")

	}
	if v.Dlvol != true {
		t.Fatalf("Dlvol flag is wrong")
	}

	if v.Ulvol != true {
		t.Fatalf("Dlvol flag is wrong")
	}

	if v.TotalVolume != 0 {
		t.Fatalf("Total volume is wrong")
	}

	if v.UplinkVolume != 0 {
		t.Fatalf("Uplink volume is wrong")
	}

	if v.DownlinkVolume != 0 {
		t.Fatalf("Downlink volume is wrong")
	}
}
