package ie

import (
	"bytes"
	"encoding/binary"
)

type SDFFilter struct {
	FD                      bool
	TTC                     bool
	SPI                     bool
	FL                      bool
	BID                     bool
	LengthOfFlowDescription uint16
	FlowDescription         []byte
	ToSTrafficClass         []byte
	SecurityParameterIndex  []byte
	FlowLabel               []byte
	SDFFilterId             uint32
}

func (s *SDFFilter) NewSDFFilterFromBytes(b []byte) error {
	buffer := bytes.NewBuffer(b)
	b1, err := buffer.ReadByte()
	if err != nil {
		return err
	}
	if (b1 & 0x01) == 0x01 {
		s.FD = true
		s.LengthOfFlowDescription = binary.BigEndian.Uint16(buffer.Next(2))
		s.FlowDescription = buffer.Next(int(s.LengthOfFlowDescription))
	}
	if (b1 & 0x02) == 0x02 {
		s.TTC = true
		s.ToSTrafficClass = buffer.Next(2)
	}
	if (b1 & 0x04) == 0x04 {
		s.SPI = true
		s.SecurityParameterIndex = buffer.Next(4)
	}
	if (b1 & 0x08) == 0x08 {
		s.FL = true
		s.FlowLabel = buffer.Next(3)
	}
	if (b1 & 0x10) == 0x10 {
		s.BID = true
		s.SDFFilterId = binary.BigEndian.Uint32(buffer.Next(4))
	}
	return nil
}

func (s *SDFFilter) Serialize() ([]byte, error) {
	b := make([]byte, 2)
	b[1] = 0x00
	if s.FD {
		b[0] = 0x01
		d := make([]byte, 2)
		binary.BigEndian.PutUint16(d, s.LengthOfFlowDescription)
		b = append(b, d...)
		b = append(b, s.FlowDescription...)
	}
	if s.TTC {
		b[0] |= 0x02
		b = append(b, s.ToSTrafficClass...)
	}
	if s.SPI {
		b[0] |= 0x04
		b = append(b, s.SecurityParameterIndex...)
	}
	if s.FL {
		b[0] |= 0x08
		b = append(b, s.FlowLabel...)
	}
	if s.BID {
		b[0] |= 0x10
		d := make([]byte, 4)
		binary.BigEndian.PutUint32(d, s.SDFFilterId)
		b = append(b, d...)
	}
	return b, nil
}
