package ie

import "encoding/binary"

type DroppedDLTrafficThreshold struct {
	DLPA                        bool
	DLBY                        bool
	DownlinkPackets             uint64
	NumberofBytesofDownlinkData uint64
}

func NewDroppedDLTrafficThreshold(dlpa, dlby bool, dp, nb uint64) *DroppedDLTrafficThreshold {
	return &DroppedDLTrafficThreshold{
		DLPA:                        dlpa,
		DLBY:                        dlby,
		DownlinkPackets:             dp,
		NumberofBytesofDownlinkData: nb,
	}

}
func (d DroppedDLTrafficThreshold) Serialize() ([]byte, error) {
	var b []byte
	var fByte uint8
	if d.DLPA {
		fByte = 0x01
	}
	if d.DLBY {
		fByte |= 0x02
	}
	b = append(b, fByte)
	if d.DLPA {

		data := make([]byte, 8)
		binary.BigEndian.PutUint64(data, d.DownlinkPackets)
		b = append(b, data...)
	}
	if d.DLBY {
		data := make([]byte, 8)
		binary.BigEndian.PutUint64(data, d.NumberofBytesofDownlinkData)
		b = append(b, data...)

	}
	return b, nil
}
