package ie

import (
	//	"bitbucket.org/sothy5/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

type CreateURR struct {
	URRID             *uint32             `tlv:"81"`
	MeasurementMethod *MeasurementMethod  `tlv:"62"`
	ReportingTriggers *ReportingTriggers  `tlv:"37"`
	MeasurementPeriod *MeasurementPeriod  `tlv:"64"`
	VolumeThreshold   *Volume             `tlv:"31"` //
	VolumeQuota       *Volume             `tlv:"73"`
	TimeThreshold     *TimeThresholdValue `tlv:"32"`
	TimeQuota         *TimeQuotaValue     `tlv:"74"`
	//QuotaHoldingTime          *QuotaHoldingTime          `tlv:"71"`
	DroppedDLTrafficThreshold *DroppedDLTrafficThreshold `tlv:"72"`
	MonitoringTime            *MonitoringTime            `tlv:"33"`
	/*
		EventInformation          *EventInformation                   `tlv:"148"`
		//SubsequentVolumeThreshold *pfcpType.SubsequentVolumeThreshold `tlv:"34"`
		SubsequentTimeThreshold   *pfcpType.SubsequentTimeThreshold   `tlv:"35"`
		SubsequentVolumeQuota     *pfcpType.SubsequentVolumeQuota     `tlv:"121"`
		SubsequentTimeQuota       *pfcpType.SubsequentTimeQuota       `tlv:"122"`
		InactivityDetectionTime   *pfcpType.InactivityDetectionTime   `tlv:"36"`
		LinkedURRID               *pfcpType.LinkedURRID               `tlv:"82"`
		MeasurementInformation    *pfcpType.MeasurementInformation    `tlv:"100"`
		TimeQuotaMechanism        *pfcpType.TimeQuotaMechanism        `tlv:"115"`
		AggregatedURRs            *AggregatedURRs                     `tlv:"118"`
		FARIDForQuotaAction       *pfcpType.FARID                     `tlv:"108"`
		EthernetInactivityTimer   *pfcpType.EthernetInactivityTimer   `tlv:"146"`
		AdditionalMonitoringTime  *AdditionalMonitoringTime           `tlv:"147"`
	*/
}

//Serialize produces the binary for CreateURR IE
func (c CreateURR) Serialize() ([]byte, error) {
	urrIE := NewInformationElement(
		IEURRID,
		0,
		dt.Unsigned32(uint32(*c.URRID)),
	)
	data, err := urrIE.Serialize()
	if err != nil {
		return nil, err
	}

	b, err := c.MeasurementMethod.Serialize()
	if err != nil {
		return nil, err
	}
	mmIE := NewInformationElement(
		IEMeasurementMethod,
		0,
		dt.OctetString(b),
	)
	data1, err := mmIE.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, data1...)

	b1, err := c.ReportingTriggers.Serialize()
	if err != nil {
		return nil, err
	}
	rrIE := NewInformationElement(
		IEReportingTriggers,
		0,
		dt.OctetString(b1),
	)
	data1, err = rrIE.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, data1...)

	if c.ReportingTriggers.PERIO && uint32(*c.MeasurementPeriod) > 0 {
		mpIE := NewInformationElement(
			IEMeasurementPeriod,
			0,
			dt.Unsigned32(uint32(*c.MeasurementPeriod)),
		)
		data1, err = mpIE.Serialize()
		if err != nil {
			return nil, err
		}
		data = append(data, data1...)
	}

	if c.ReportingTriggers.VOLTH && c.MeasurementMethod.VOLUM {
		b1, err = c.VolumeThreshold.Serilize()
		if err != nil {
			return nil, err
		}
		vtIE := NewInformationElement(
			IEVolumeThreshold,
			0,
			dt.OctetString(b1),
		)
		data1, err = vtIE.Serialize()
		if err != nil {
			return nil, err
		}
		data = append(data, data1...)
	}
	if c.ReportingTriggers.VOLQU && c.MeasurementMethod.VOLUM {
		b1, err = c.VolumeThreshold.Serilize()
		if err != nil {
			return nil, err
		}
		vtIE := NewInformationElement(
			IEVolumeQuota,
			0,
			dt.OctetString(b1),
		)
		data1, err = vtIE.Serialize()
		if err != nil {
			return nil, err
		}
		data = append(data, data1...)
	}

	if c.ReportingTriggers.TIMTH && c.MeasurementMethod.DURAT {
		ttIE := NewInformationElement(
			IETimeThreshold,
			0,
			dt.Unsigned32(uint32(*c.TimeThreshold)),
		)
		data1, err = ttIE.Serialize()
		if err != nil {
			return nil, err
		}
		data = append(data, data1...)
	}

	if c.ReportingTriggers.TIMQU && c.MeasurementMethod.DURAT {
		ttIE := NewInformationElement(
			IETimeQuota,
			0,
			dt.Unsigned32(uint32(*c.TimeQuota)),
		)
		data1, err = ttIE.Serialize()
		if err != nil {
			return nil, err
		}
		data = append(data, data1...)
	}

	return data, nil

}
