package ie

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
