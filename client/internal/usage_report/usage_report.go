package usage_report

import (
	"bitbucket.org/sothy5/n4-go/ie"
	"bitbucket.org/sothy5/n4-go/ie/urr"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

func NewCreateURR(urrid uint32, m *urr.MeasurementMethod, r *urr.ReportingTriggers, mp uint32, tt uint32) (*ie.InformationElement, error) {
	//TODO Time Quota/Quota Holding Time/Subsequent Time Threshold/Inactivity Detection Time Value is duration <<uint32>>
	urrIE := ie.NewInformationElement(
		ie.IEURRID,
		0,
		dt.Unsigned32(urrid),
	)
	data, err := urrIE.Serialize()
	if err != nil {
		return nil, err
	}

	b, err := m.Serialize()
	if err != nil {
		return nil, err
	}
	mmIE := ie.NewInformationElement(
		ie.IEMeasurementMethod,
		0,
		dt.OctetString(b),
	)
	data1, err := mmIE.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, data1...)

	b, err = m.Serialize()
	if err != nil {
		return nil, err
	}
	rrIE := ie.NewInformationElement(
		ie.IEReportingTriggers,
		0,
		dt.OctetString(b),
	)
	data1, err = rrIE.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, data1...)

	if r.PERIO && mp > 0 {
		mpIE := ie.NewInformationElement(
			ie.IEMetric,
			0,
			dt.Unsigned32(mp),
		)
		data1, err = mpIE.Serialize()
		if err != nil {
			return nil, err
		}
		data = append(data, data1...)
	}

	ttIE := ie.NewInformationElement(
		ie.IETimeThreshold,
		0,
		dt.Unsigned32(tt),
	)
	data1, err = ttIE.Serialize()
	if err != nil {
		return nil, err
	}
	data = append(data, data1...)

	createurr := ie.NewInformationElement(
		ie.IECreateURR,
		0,
		dt.OctetString(data),
	)

	return &createurr, nil

}
