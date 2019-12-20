package sr

import (
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
