package ie

import "fmt"

type MeasurementMethod struct {
	DURAT bool
	VOLUM bool
	EVENT bool
}

func NewMeasurementMethod(d, v, e bool) *MeasurementMethod {
	return &MeasurementMethod{
		DURAT: d,
		VOLUM: v,
		EVENT: e,
	}
}

func (m *MeasurementMethod) Serialize() (byte, error) {
	if m.DURAT == false && m.VOLUM == false && m.EVENT == false {
		return 0, fmt.Errorf("at least one bit shall be set to 1")
	}
	var b uint8
	if m.DURAT {
		b = 1
	}
	if m.VOLUM {
		b |= 0x02
	}
	if m.EVENT {
		b |= 0x04
	}
	return b, nil
}
