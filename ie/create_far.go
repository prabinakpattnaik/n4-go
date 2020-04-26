package ie

import (
	"fmt"

	dt "github.com/fiorix/go-diameter/diam/datatype"
)

type CreateFAR struct {
	FarID                *FARID                       `tlv:"108"`
	ApplyAction          *ApplyAction                 `tlv:"44"`
	ForwardingParameters *ForwardingParametersIEInFAR `tlv:"4"`
	//DuplicatingParameters *DuplicatingParameters          `tlv:"5"`
	//BARID                 *BARID                          `tlv:"88"`
}

func (c *CreateFAR) Serialize() ([]byte, error) {
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
	if c.ForwardingParameters != nil {
		bb, err := c.ForwardingParameters.Serialize()
		if err != nil {
			return nil, err
		}

		fpIE = NewInformationElement(
			IEForwardingParameters,
			0,
			dt.OctetString(bb),
		)
	}
	createFAR := NewCreateFAR(&f, &a, &fpIE, nil, nil)
	return createFAR.Serialize()

}
