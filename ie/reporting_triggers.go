package ie

type ReportingTriggers struct {
	PERIO bool
	VOLTH bool
	TIMTH bool
	QUHTI bool
	START bool
	STOPT bool
	DROTH bool
	LIUSA bool
	VOLQU bool
	TIMQU bool
	ENVCL bool
	MACAR bool
	EVETH bool
	EVEQU bool
}

func NewReportingTriggers(p, v, t, q, start, stopt, d, l, vol, tim, en, ma, eveth, evequ bool) *ReportingTriggers {
	return &ReportingTriggers{
		PERIO: p,
		VOLTH: v,
		TIMTH: t,
		QUHTI: q,
		START: start,
		STOPT: stopt,
		DROTH: d,
		LIUSA: l,
		VOLQU: vol,
		TIMQU: tim,
		ENVCL: en,
		MACAR: ma,
		EVETH: eveth,
		EVEQU: evequ,
	}
}

func (r ReportingTriggers) Serialize() ([]byte, error) {
	var fByte, sByte uint8
	if r.PERIO {
		fByte = 1
	}
	if r.VOLTH {
		fByte |= 0x02
	}
	if r.TIMTH {
		fByte |= 0x04
	}
	if r.QUHTI {
		fByte |= 0x08
	}
	if r.START {
		fByte |= 0x10
	}
	if r.STOPT {
		fByte |= 0x20
	}
	if r.DROTH {
		fByte |= 0x40
	}
	if r.LIUSA {
		fByte |= 0x80
	}
	if r.VOLQU {
		sByte = 0x01
	}
	if r.TIMQU {
		sByte |= 0x02
	}
	if r.ENVCL {
		sByte |= 0x04
	}
	if r.MACAR {
		sByte |= 0x08
	}
	if r.EVETH {
		sByte |= 0x10
	}
	if r.EVEQU {
		sByte |= 0x20
	}
	b := make([]byte, 2)
	b[0] = fByte
	b[1] = sByte
	return b, nil

}
