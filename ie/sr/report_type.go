package sr

// ReportType
type ReportType struct {
	DLDR bool
	USAR bool
	ERIR bool
	UPIR bool
}

func NewReportType(d, u, e, up bool) *ReportType {
	return &ReportType{
		DLDR: d,
		USAR: u,
		ERIR: e,
		UPIR: up,
	}

}

func (r ReportType) Serialize() byte {
	var b uint8
	if r.DLDR {
		b = 1
	}
	if r.USAR {
		b = b | 0x02
	}

	if r.ERIR {
		b = b | 0x04
	}
	if r.UPIR {
		b = b | 0x08
	}

	return b
}

func NewReportTypeFromByte(b byte) *ReportType {
	d := ((b & 0x01) == 1)
	u := ((b & 0x02) == 2)
	e := ((b & 0x04) == 4)
	up := ((b & 0x08) == 8)
	return NewReportType(d, u, e, up)

}
