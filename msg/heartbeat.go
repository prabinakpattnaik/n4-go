package msg

import (
	"bitbucket.org/sothy5/n4-go/ie"
)

type Heartbeat struct {
	Header            *PFCPHeader
	RecoveryTimeStamp *ie.InformationElement
}

func NewHeartbeat(h *PFCPHeader, r *ie.InformationElement) *Heartbeat {
	return &Heartbeat{
		Header:            h,
		RecoveryTimeStamp: r,
	}

}

func (hr Heartbeat) Serialize() ([]byte, error) {
	b, err := hr.RecoveryTimeStamp.Serialize()
	if err != nil {
		return nil, err
	}
	output := make([]byte, hr.Len())
	pfcpend := uint16(PFCPBasicHeaderLength) + PFCPBasicMessageSize
	copy(output[:pfcpend], hr.Header.Serialize())
	copy(output[pfcpend:], b)

	return output, nil

}

func (hr Heartbeat) Len() uint16 {

	return uint16(PFCPBasicHeaderLength) + hr.Header.MessageLength

}
