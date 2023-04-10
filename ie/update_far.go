package ie

import (
	"fmt"

	dt "github.com/fiorix/go-diameter/diam/datatype"
)

type UpdateFAR struct {
	FarID                      *FARID                       `tlv:"108"`
	ApplyAction                *ApplyAction                 `tlv:"44"`
	UpdateForwardingParameters *UpdateForwardingParametersIEInFAR `tlv:"11"`
	//UpdateDuplicatingParameters *pfcpType.UpdateDuplicatingParameters `tlv:"105"`
	//BARID                 *BARID                          `tlv:"88"`
}

func (c *UpdateFAR) Serialize() ([]byte, error) {
	if c.FarID == nil {
		return nil, fmt.Errorf("Mandory field is missing in Create FAR(FARID)")
	}
	f := NewInformationElement(
		IEFARID,
		0,
		dt.Unsigned32(c.FarID.FarIdValue),
	)

	if c.ApplyAction == nil {
		return nil, fmt.Errorf("Mandory field is missing in Create FAR(ApplyAction)")
	}
	b, err := c.ApplyAction.Serialize()
	if err != nil {
		return nil, err
	}
	applyAction := dt.OctetString(b)
	a := NewInformationElement(
		IEApplyAction,
		0,
		applyAction,
	)
	var fpIE InformationElement
	if c.UpdateForwardingParameters != nil {
		bb, err := c.UpdateForwardingParameters.Serialize()
		if err != nil {
			return nil, err
		}

		fpIE = NewInformationElement(
			IEUpdateForwardingParamets,
			0,
			dt.OctetString(bb),
		)
	}
	updateFAR := NewUpdateFAR(&f, &a, &fpIE, nil, nil)
	return updateFAR.Serialize()

}
