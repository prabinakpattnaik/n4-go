package sr

import (
	"encoding/binary"
	"time"

	"github.com/prabinakpattnaik/n4-go/ie"
)

// UsageReport
type UsageReport struct {
	URRID               *ie.URRID
	URSEQN              *ie.URSEQN
	UsageReportTrigger  *ie.UsageReportTriggerData
	StartTime           *time.Time
	EndTime             *time.Time
	VolumeMeasurement   *ie.Volume
	DurationMeasurement *uint32
	//ApplicationDetectionInformation *ie.InformationElement
	UEIPAddress       *ie.UEIPAddress
	NetworkInstance   []byte
	TimeofFirstPacket *time.Time
	TimeofLastPacket  *time.Time
	UsageInformation  []byte
	//QueryURRReference               *ie.InformationElement
	//EventTimeStamp                  *ie.InformationElement
	//EthernetTrafficInformation      *ie.InformationElement
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
			urrID := binary.BigEndian.Uint32(inel.Data.Serialize())
			uid := ie.URRID(urrID)
			u.URRID = &uid
		case ie.IEUESEQN:
			urseqn := binary.BigEndian.Uint32(inel.Data.Serialize())
			urs := ie.URSEQN(urseqn)
			u.URSEQN = &urs
		case ie.IEUsageReportTrigger:
			ur := ie.UsageReportTriggerData(inel.Data.Serialize())
			u.UsageReportTrigger = &ur
		case ie.IEStartTime:
			b := inel.Data.Serialize()
			if len(b) == 4 {
				st := time.Unix(int64(binary.BigEndian.Uint32(b))-ie.RFC868offset, 0)
				u.StartTime = &st
			}
		case ie.IEEndTime:
			b := inel.Data.Serialize()
			if len(b) == 4 {
				et := time.Unix(int64(binary.BigEndian.Uint32(b))-ie.RFC868offset, 0)
				u.EndTime = &et
			}
		case ie.IEVolumeMeasurement:
			var vol ie.Volume
			err := vol.VolumeFromBytes(inel.Data.Serialize())
			if err == nil {
				u.VolumeMeasurement = &vol
			}
		case ie.IEDurationMeasurement:
			dm := binary.BigEndian.Uint32(b)
			u.DurationMeasurement = &dm
		default:
			//fmt.Errorf("IEs are not treated")

		}

	}
	return nil

}
