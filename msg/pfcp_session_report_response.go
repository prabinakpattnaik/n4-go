package msg

import (
	"github.com/prabinakpattnaik/n4-go/ie"
)

//PFCPSessionReportResponse
type PFCPSessionReportResponse struct {
	Header         *PFCPHeader
	Cause          *ie.InformationElement
	OffendingIE    *ie.InformationElement
	UpdateBAR      *ie.InformationElement
	PFCPSRRspFlags *ie.InformationElement
}

func NewPFCPSessionReportResponse(h *PFCPHeader, c, o, u, p *ie.InformationElement) PFCPSessionReportResponse {
	return PFCPSessionReportResponse{
		Header:         h,
		Cause:          c,
		OffendingIE:    o,
		UpdateBAR:      u,
		PFCPSRRspFlags: p,
	}
}

func (sr PFCPSessionReportResponse) Serialize() ([]byte, error) {
	dataLength := sr.Len()
	output := make([]byte, dataLength)
	pfcpend := uint16(PFCPBasicHeaderLength) + PFCPMessageSize
	copy(output[:pfcpend], sr.Header.Serialize())

	b, err := sr.Cause.Serialize()
	if err != nil {
		return nil, err
	}

	if sr.OffendingIE != nil {
		b1, err := sr.OffendingIE.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, b1...)
	}

	if sr.UpdateBAR != nil {
		b1, err := sr.UpdateBAR.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, b1...)

	}

	copy(output[pfcpend:], b)
	return output, nil
}

func (sr PFCPSessionReportResponse) Len() uint16 {
	return uint16(PFCPBasicHeaderLength) + sr.Header.MessageLength
}

func (sr PFCPSessionReportResponse) Type() PFCPType {
	return sr.Header.MessageType
}

func (sr PFCPSessionReportResponse) GetHeader() *PFCPHeader {
	return sr.Header
}
