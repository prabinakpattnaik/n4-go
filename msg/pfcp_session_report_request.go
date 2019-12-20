package msg

import (
	"bitbucket.org/sothy5/n4-go/ie"
)

//PFCPSessionReportRequest
type PFCPSessionReportRequest struct {
	Header                            *PFCPHeader
	ReportType                        *ie.InformationElement
	DownlinkDataReport                *ie.InformationElement
	UsageReport                       ie.InformationElements
	ErrorIndicationReport             *ie.InformationElement
	LoadControlInformation            *ie.InformationElement
	OverloadControlInformation        *ie.InformationElement
	AdditionalUsageReportsInformation *ie.InformationElement
}

func NewPFCPSessionReportRequest(h *PFCPHeader, r, d, e, l, o, a *ie.InformationElement, u ie.InformationElements) PFCPSessionReportRequest {
	return PFCPSessionReportRequest{
		Header:                            h,
		ReportType:                        r,
		DownlinkDataReport:                d,
		UsageReport:                       u,
		ErrorIndicationReport:             e,
		LoadControlInformation:            l,
		OverloadControlInformation:        o,
		AdditionalUsageReportsInformation: a,
	}
}

func (sr PFCPSessionReportRequest) Serialize() ([]byte, error) {
	dataLength := sr.Len()
	output := make([]byte, dataLength)
	pfcpend := uint16(PFCPBasicHeaderLength) + PFCPMessageSize
	copy(output[:pfcpend], sr.Header.Serialize())

	b, err := sr.ReportType.Serialize()
	if err != nil {
		return nil, err
	}

	if sr.DownlinkDataReport != nil {
		b1, err := sr.DownlinkDataReport.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, b1...)
	}

	copy(output[pfcpend:], b)
	return output, nil

}

func (sr PFCPSessionReportRequest) Len() uint16 {
	return uint16(PFCPBasicHeaderLength) + sr.Header.MessageLength
}

func (sr PFCPSessionReportRequest) Type() PFCPType {
	return sr.Header.MessageType
}

func (sr PFCPSessionReportRequest) GetHeader() *PFCPHeader {
	return sr.Header
}
