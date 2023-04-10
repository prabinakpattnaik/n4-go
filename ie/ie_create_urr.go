package ie

import (
        "fmt"
)

//CreateURR is a struct of createURR
type CreateURRWithIE struct {
	URRID                     *InformationElement
	MeasurementMethod         *InformationElement
	ReportingTriggers         *InformationElement
	MeasurementPeriod         *InformationElement
	VolumeThreshold           *InformationElement
	VolumeQuota               *InformationElement
	EventThreshold            *InformationElement
	EventQuota                *InformationElement
	TimeThreshold             *InformationElement
	TimeQuota                 *InformationElement
	QuotaHoldingTime          *InformationElement
	DroppedDLTrafficThreshold *InformationElement
	MonitoringTime            *InformationElement
	SubsequentVolumeThreshold *InformationElement
	SubsequentTimeThreshold   *InformationElement
	SubsequentVolumeQuota     *InformationElement
	SubsequentTimeQuota       *InformationElement
	SubsequentEventThreshold  *InformationElement
	SubsequentEventQuota      *InformationElement
	InactivityDetectionTime   *InformationElement
	LinkedURRID               *InformationElement
	MeasurementInformation    *InformationElement
	FARIDForQuotaAction       *InformationElement
	EthernetInactivityTimer   *InformationElement
	AdditionalMonitoringTime  *InformationElement
}

//NewCreateURR creates new CreateURR struct
func NewCreateURR(u, mm, r, mp, vt, vq, et, eq, tt, tq, qht, d, mt *InformationElement) *CreateURRWithIE {
	return &CreateURRWithIE{
		URRID:                     u,
		MeasurementMethod:         mm,
		ReportingTriggers:         r,
		MeasurementPeriod:         mp,
		VolumeThreshold:           vt,
		VolumeQuota:               vq,
		EventThreshold:            et,
		EventQuota:                eq,
		TimeThreshold:             tt,
		TimeQuota:                 tq,
		QuotaHoldingTime:          qht,
		DroppedDLTrafficThreshold: d,
		MonitoringTime:            mt,
		SubsequentVolumeThreshold: nil,
		SubsequentTimeThreshold:   nil,
		SubsequentVolumeQuota:     nil,
		SubsequentTimeQuota:       nil,
		SubsequentEventThreshold:  nil,
		SubsequentEventQuota:      nil,
		InactivityDetectionTime:   nil,
		LinkedURRID:               nil,
		MeasurementInformation:    nil,
		FARIDForQuotaAction:       nil,
		EthernetInactivityTimer:   nil,
		AdditionalMonitoringTime:  nil,
	}

}
func (c CreateURRWithIE) Serialize() ([]byte, error) {
        if c.URRID == nil || c.URRID.Type == IEReserved {
                return nil, fmt.Errorf("CreateFAR does not have valid URRID")
        }
        b, err := c.URRID.Serialize()
        if err != nil {
                return nil, fmt.Errorf("URRID serialization error")
        }
        if c.MeasurementMethod == nil || c.MeasurementMethod.Type == IEReserved {
                return nil, fmt.Errorf("CreateFAR does not have valid MeasurementMethod")
        }
        b1, err := c.MeasurementMethod.Serialize()
        if err != nil {
                return nil, fmt.Errorf("MeasurementMethod serialization error")
        }
        b = append(b, b1...)

        if c.ReportingTriggers != nil && c.ReportingTriggers.Type != IEReserved {
		b1, err := c.ReportingTriggers.Serialize()
                if err != nil {
			return nil, fmt.Errorf("ReportingTriggers serialization error")
		}
                b = append(b, b1...)
	}
	if c.MeasurementPeriod != nil && c.MeasurementPeriod.Type != IEReserved {
                b1, err := c.MeasurementPeriod.Serialize()
                if err != nil {
			return nil, fmt.Errorf("MeasurementPeriod serialization error")
                }
                b = append(b, b1...)
	}
	if c.VolumeThreshold != nil && c.VolumeThreshold.Type != IEReserved {
                b1, err := c.VolumeThreshold.Serialize()
                if err != nil {
                        return nil, fmt.Errorf("VolumeThreshold serialization error")
                }
                b = append(b, b1...)
	}
	if c.VolumeQuota != nil && c.VolumeQuota.Type != IEReserved {
                b1, err := c.VolumeQuota.Serialize()
                if err != nil {
                        return nil, fmt.Errorf("VolumeQuota serialization error")
                }
                b = append(b, b1...)
	}
	if c.EventThreshold != nil && c.EventThreshold.Type != IEReserved {
                b1, err := c.EventThreshold.Serialize()
                if err != nil {
                        return nil, fmt.Errorf("EventThreshold serialization error")
		}
                b = append(b, b1...)
	}
	if c.EventQuota != nil && c.EventQuota.Type != IEReserved {
                b1, err := c.EventQuota.Serialize()
                if err != nil {
                        return nil, fmt.Errorf("EventQuota serialization error")
                }
                b = append(b, b1...)
        }
	if c.TimeThreshold != nil && c.TimeThreshold.Type != IEReserved {
                b1, err := c.TimeThreshold.Serialize()
                if err != nil {
                        return nil, fmt.Errorf("TimeThreshold serialization error")
                }
                b = append(b, b1...)
        }

        return b, nil
}
