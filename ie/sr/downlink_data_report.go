package sr

import (
	"github.com/prabinakpattnaik/n4-go/ie"
	 "encoding/binary"
 )

type DownlinkDataReport struct {
	PDRID                            uint16
	DownlinkDataServiceInformation   *DownlinkDataServiceInformation
}

func  NewDownlinkDataReportFromBytes ( b [] byte) *DownlinkDataReport {
     var ddrIEs ie.InformationElements
     err := ddrIEs.FromBytes(b)
	if err != nil {
		return nil
	}
	var ddr DownlinkDataReport

	for  _, informationElement := range ddrIEs {
	   switch informationElement.Type {
	   case ie.IEPDRID:
		   ddr.PDRID=binary.BigEndian.Uint16(informationElement.Data.Serialize())
           case ie.IEDownlinkDataServiceInformation:
		   ddsi, err:=NewDownlinkDataServiceInformationFromBytes(informationElement.Data.Serialize())
		   if err !=nil {
		   ddr.DownlinkDataServiceInformation=ddsi
	           }
	   }
	}
	return &ddr

}
