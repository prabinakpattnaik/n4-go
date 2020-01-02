package bar

import (
	"fmt"

	"bitbucket.org/sothy5/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

//CreateBAR is a struct of createBAR
type CreateBAR struct {
	BARID                          *ie.InformationElement
	SuggestedBufferingPacketsCount *ie.InformationElement
}

// value between 0-255
type PacketCountvalue uint8

func NewCreateBAR(barid uint8, pcv PacketCountvalue) (*CreateBAR, error) {
	barIDIE := ie.NewInformationElement(
		ie.IEBARID,
		0,
		dt.OctetString(barid),
	)
	sbpcIE := ie.NewInformationElement(
		ie.IESuggestedBufferingPacketsCount,
		0,
		dt.OctetString(pcv),
	)

	createBAR := CreateBAR{
		BARID:                          &barIDIE,
		SuggestedBufferingPacketsCount: &sbpcIE,
	}
	return &createBAR, nil
}

func (c *CreateBAR) Serialize() ([]byte, error) {
	if c.BARID == nil {
		return nil, fmt.Errorf("BARID is nil")
	}
	if c.BARID.Type == ie.IEReserved {
		return nil, fmt.Errorf("BARID does not have valid Type")
	}
	b, err := c.BARID.Serialize()
	if err != nil {
		return nil, fmt.Errorf("CreateBAR: BARID serialization error")
	}
	if c.SuggestedBufferingPacketsCount != nil && c.SuggestedBufferingPacketsCount.Type != ie.IEReserved {
		b1, err := c.SuggestedBufferingPacketsCount.Serialize()
		if err != nil {
			return nil, fmt.Errorf("CreateBAR: SuggestedBufferingPacketsCount serialization error")
		}
		b = append(b, b1...)
	}
	return b, nil

}
