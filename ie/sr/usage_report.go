package sr

import (
	"fmt"

	"bitbucket.org/sothy5/n4-go/ie"
)

// DownlinkDataReport
type UsageReport struct {
	URRID                           *ie.InformationElement
	URSEQN                          *ie.InformationElement
	UsageReportTrigger              *ie.InformationElement
	StartTime                       *ie.InformationElement
	EndTime                         *ie.InformationElement
	VolumeMeasurement               *ie.InformationElement
	DurationMeasurement             *ie.InformationElement
	ApplicationDetectionInformation *ie.InformationElement
	UEIPAddress                     *ie.InformationElement
	NetworkInstance                 *ie.InformationElement
	TimeofFirstPacket               *ie.InformationElement
	TimeofLastPacket                *ie.InformationElement
	UsageInformation                *ie.InformationElement
	QueryURRReference               *ie.InformationElement
	EventTimeStamp                  *ie.InformationElement
	EthernetTrafficInformation      *ie.InformationElement
}

func (u *UsageReport) UsageReportFromByte(b []byte) error {
	var ies ie.InformationElements
	err := ies.FromBytes(b)
	if err != nil {
		return err
	}
	for _, inel := range ies {
		switch inel.Type {
		case ie.IEURRID:
			u.URRID = &inel
		case ie.IEUESEQN:
			u.URSEQN = &inel
		case ie.IEUsageReportTrigger:
			u.UsageReportTrigger = &inel
		case ie.IEStartTime:
			u.StartTime = &inel
		case ie.IEEndTime:
			u.EndTime = &inel
		case ie.IEVolumeMeasurement:
			u.VolumeMeasurement = &inel
		case ie.IEDurationMeasurement:
			u.DurationMeasurement = &inel
		default:
			return fmt.Errorf("IEs are not treated")

		}

	}
	return nil

}
