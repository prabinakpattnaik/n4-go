package sr

// DownlinkDataReport
type DownlinkDataServiceInformation struct {
	PPI      bool
	QFII     bool
	PPIValue uint8
	QFI      uint8
}

func NewDownlinkDataServiceInformationFromBytes(b []byte)(*DownlinkDataServiceInformation,error) {
     firstbyte :=b[0]
     var ddsi DownlinkDataServiceInformation
     if (firstbyte &0x01)==0x01{
       ddsi.PPI=true
       ddsi.PPIValue= b[1]
     }
     if (firstbyte & 0x02)==0x02 {
       ddsi.QFII=true
       if ddsi.PPI==true {
       ddsi.QFI=b[2]
       }else  {
       ddsi.QFI=b[1]
       }
     }
    return &ddsi, nil
  
}
