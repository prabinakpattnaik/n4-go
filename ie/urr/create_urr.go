package urr

import "bitbucket.org/sothy5/n4-go/ie"

//CreateURR is a struct of createURR
type CreateURR struct {
	URRID                     *ie.InformationElement
	MeasurementMethod         *ie.InformationElement
	ReportingTriggers         *ie.InformationElement
	MeasurementPeriod         *ie.InformationElement
	VolumeThreshold           *ie.InformationElement
	VolumeQuota               *ie.InformationElement
	EventThreshold            *ie.InformationElement
	EventQuota                *ie.InformationElement
	TimeThreshold             *ie.InformationElement
	TimeQuota                 *ie.InformationElement
	QuotaHoldingTime          *ie.InformationElement
	DroppedDLTrafficThreshold *ie.InformationElement
	MonitoringTime            *ie.InformationElement
	SubsequentVolumeThreshold *ie.InformationElement
	SubsequentTimeThreshold   *ie.InformationElement
	SubsequentVolumeQuota     *ie.InformationElement
	SubsequentTimeQuota       *ie.InformationElement
	SubsequentEventThreshold  *ie.InformationElement
	SubsequentEventQuota      *ie.InformationElement
	InactivityDetectionTime   *ie.InformationElement
	LinkedURRID               *ie.InformationElement
	MeasurementInformation    *ie.InformationElement
	FARIDForQuotaAction       *ie.InformationElement
	EthernetInactivityTimer   *ie.InformationElement
	AdditionalMonitoringTime  *ie.InformationElement
}

//NewCreateURR creates new CreateURR struct
func NewCreateURR(u, mm, r, mp, vt, vq, et, eq, tt, tq, qht, d, mt *ie.InformationElement) *CreateURR {
	return &CreateURR{
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
