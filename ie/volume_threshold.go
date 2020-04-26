package ie

import (
	"encoding/binary"
	"fmt"
)

type Volume struct {
	Dlvol          bool
	Ulvol          bool
	Tovol          bool
	TotalVolume    uint64
	UplinkVolume   uint64
	DownlinkVolume uint64
}

func (v Volume) Serilize() ([]byte, error) {
	b := make([]byte, 1)
	if v.Tovol {
		b[0] = 1
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v.TotalVolume)
		b = append(b, buf...)
	}
	if v.Ulvol {
		b[0] |= 2
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v.UplinkVolume)
		b = append(b, buf...)
	}
	if v.Dlvol {
		b[0] |= 4
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v.DownlinkVolume)
		b = append(b, buf...)
	}
	return b, nil

}

func (v *Volume) VolumeFromBytes(b []byte) error {
	if (b[0] & 0x07) == 0x00 {
		return fmt.Errorf("No flag set for volume")
	}

	if uint8(b[0]&0x01) == 1 {
		v.Tovol = true
		v.TotalVolume = binary.BigEndian.Uint64(b[1:9])
	}
	if (b[0] & 0x02) == 2 {
		v.Ulvol = true
		if v.Tovol == true {
			v.UplinkVolume = binary.BigEndian.Uint64(b[9:17])
		} else {
			v.UplinkVolume = binary.BigEndian.Uint64(b[1:9])

		}
	}

	if (b[0] & 0x04) == 0x04 {
		v.Dlvol = true
		if v.Tovol && v.Dlvol {
			v.DownlinkVolume = binary.BigEndian.Uint64(b[17:])
		} else if (v.Tovol && !v.Dlvol) || (!v.Tovol && v.Dlvol) {
			v.DownlinkVolume = binary.BigEndian.Uint64(b[9:17])
		} else {
			v.DownlinkVolume = binary.BigEndian.Uint64(b[1:9])
		}

	}
	return nil

}
