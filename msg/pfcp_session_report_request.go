package msg

import (
	"fmt"

	"bitbucket.org/sothy5/n4-go/ie"
	"bitbucket.org/sothy5/n4-go/ie/sr"
	"bitbucket.org/sothy5/n4-go/util/se"
	dt "github.com/fiorix/go-diameter/diam/datatype"
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
func ProcessPFCPSessionReportRequest(m *PFCPMessage, cpSEIDDPSEID *se.CPSEIDDPSEIDEntity) ([]byte, error) {
	pfcpMessage, err := FromPFCPMessage(m)
	if err != nil {
		fmt.Printf("pfcpSessionER:Header %+v\n", pfcpMessage)
		return nil, err
	}
	//TODO bad length
	pfcpSessionReportRequest, ok := pfcpMessage.(PFCPSessionReportRequest)
	if !ok {
		return nil, fmt.Errorf("PFCP Session Report Request could not type asserted")
	}
	pfcpHeader := pfcpSessionReportRequest.GetHeader()
	seid := cpSEIDDPSEID.Value(pfcpHeader.SessionEndpointIdentifier)
	pfcpHeader.MessageType = SessionReportResponseType
	if seid == 0 {
		pfcpHeader.SessionEndpointIdentifier = 0
		c := ie.NewInformationElement(
			ie.IECause,
			0,
			dt.OctetString([]byte{uint8(ie.SessionContextNotFound)}),
		)
		pfcpHeader.MessageLength = PFCPMessageSize + ie.IEBasicHeaderSize + c.Len()
		pfcpSessionReportResponse := NewPFCPSessionReportResponse(pfcpHeader, &c, nil, nil, nil)
		return pfcpSessionReportResponse.Serialize()
	}
	pfcpHeader.SessionEndpointIdentifier = seid
	if pfcpSessionReportRequest.ReportType == nil || pfcpSessionReportRequest.ReportType.Type == ie.IEReserved {
		c := ie.NewInformationElement(
			ie.IECause,
			0,
			dt.OctetString([]byte{uint8(ie.MandatoryIEMissing)}),
		)
		o := ie.NewInformationElement(
			ie.IEOffendingIE,
			0,
			dt.OctetString([]byte{uint8(pfcpSessionReportRequest.ReportType.Type)}),
		)
		pfcpHeader.MessageLength = PFCPMessageSize + ie.IEBasicHeaderSize + c.Len() + ie.IEBasicHeaderSize + o.Len()
		pfcpSessionReportResponse := NewPFCPSessionReportResponse(pfcpHeader, &c, &o, nil, nil)
		return pfcpSessionReportResponse.Serialize()
	}
	reportType := sr.NewReportTypeFromByte(pfcpSessionReportRequest.ReportType.Data.Serialize()[0])
	if reportType.DLDR {
		if pfcpSessionReportRequest.DownlinkDataReport == nil || pfcpSessionReportRequest.DownlinkDataReport.Type == ie.IEReserved {
			c := ie.NewInformationElement(
				ie.IECause,
				0,
				dt.OctetString([]byte{uint8(ie.ConditionalIEMissing)}),
			)
			o := ie.NewInformationElement(
				ie.IEOffendingIE,
				0,
				dt.OctetString([]byte{uint8(pfcpSessionReportRequest.DownlinkDataReport.Type)}),
			)
			pfcpHeader.MessageLength = PFCPMessageSize + ie.IEBasicHeaderSize + c.Len() + ie.IEBasicHeaderSize + o.Len()
			pfcpSessionReportResponse := NewPFCPSessionReportResponse(pfcpHeader, &c, &o, nil, nil)
			return pfcpSessionReportResponse.Serialize()

		} else {
			//TODO Downlink Data Report
			c := ie.NewInformationElement(
				ie.IECause,
				0,
				dt.OctetString([]byte{uint8(ie.RequestAccepted)}),
			)
			pfcpHeader.MessageLength = PFCPMessageSize + ie.IEBasicHeaderSize + c.Len()
			pfcpSessionReportResponse := NewPFCPSessionReportResponse(pfcpHeader, &c, nil, nil, nil)
			return pfcpSessionReportResponse.Serialize()

		}

	} else if reportType.USAR {
		//TODO Multi IEs needed!
		if len(pfcpSessionReportRequest.UsageReport) > 0 && pfcpSessionReportRequest.UsageReport[0].Type != ie.IEReserved {
			//TODO Usage  Report
			c := ie.NewInformationElement(
				ie.IECause,
				0,
				dt.OctetString([]byte{uint8(ie.RequestAccepted)}),
			)
			pfcpHeader.MessageLength = PFCPMessageSize + ie.IEBasicHeaderSize + c.Len()
			pfcpSessionReportResponse := NewPFCPSessionReportResponse(pfcpHeader, &c, nil, nil, nil)
			return pfcpSessionReportResponse.Serialize()

		} else {
			c := ie.NewInformationElement(
				ie.IECause,
				0,
				dt.OctetString([]byte{uint8(ie.ConditionalIEMissing)}),
			)
			o := ie.NewInformationElement(
				ie.IEOffendingIE,
				0,
				dt.OctetString([]byte{uint8(pfcpSessionReportRequest.UsageReport[0].Type)}),
			)
			pfcpHeader.MessageLength = PFCPMessageSize + ie.IEBasicHeaderSize + c.Len() + ie.IEBasicHeaderSize + o.Len()
			pfcpSessionReportResponse := NewPFCPSessionReportResponse(pfcpHeader, &c, &o, nil, nil)
			return pfcpSessionReportResponse.Serialize()

		}
	} else if reportType.ERIR {
		//TODO Multi IEs needed!
		if pfcpSessionReportRequest.ErrorIndicationReport == nil || pfcpSessionReportRequest.ErrorIndicationReport.Type == ie.IEReserved {
			c := ie.NewInformationElement(
				ie.IECause,
				0,
				dt.OctetString([]byte{uint8(ie.ConditionalIEMissing)}),
			)
			o := ie.NewInformationElement(
				ie.IEOffendingIE,
				0,
				dt.OctetString([]byte{uint8(pfcpSessionReportRequest.ErrorIndicationReport.Type)}),
			)
			pfcpHeader.MessageLength = PFCPMessageSize + ie.IEBasicHeaderSize + c.Len() + ie.IEBasicHeaderSize + o.Len()
			pfcpSessionReportResponse := NewPFCPSessionReportResponse(pfcpHeader, &c, &o, nil, nil)
			return pfcpSessionReportResponse.Serialize()

		} else {
			//TODO Usage  Report
			c := ie.NewInformationElement(
				ie.IECause,
				0,
				dt.OctetString([]byte{uint8(ie.RequestAccepted)}),
			)
			pfcpHeader.MessageLength = PFCPMessageSize + ie.IEBasicHeaderSize + c.Len()
			pfcpSessionReportResponse := NewPFCPSessionReportResponse(pfcpHeader, &c, nil, nil, nil)
			return pfcpSessionReportResponse.Serialize()

		}

	}
	return nil, fmt.Errorf("Error in creating PFCP Session Report Response")
	//pfc
	//NewReportTypeFromByte(b byte) ReportType

}
