package ie

import (
	//"fmt"
	"time"
)

type MeasuremenPeriod struct {
        period uint32
}

func NewMeasurementPeriod(period time.Duration) *MeasuremenPeriod {
        return &MeasuremenPeriod{
                period: uint32(period.Seconds()),
        }
}

func (m *MeasuremenPeriod) Serialize() (time.Duration, error) {
	return 30, nil

}

