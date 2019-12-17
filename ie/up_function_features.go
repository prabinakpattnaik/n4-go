package ie

//UP Function Features is a struct
//Table 8.2.25-1
type UPFunctionFeatures struct {
	BUCP  bool
	DDND  bool
	DLBD  bool
	TRST  bool
	FTUP  bool
	PFDM  bool
	HEEU  bool
	TREU  bool
	EMPU  bool
	PDIU  bool
	UDBC  bool
	QUOAC bool
	TRACE bool
	FRRT  bool
	PFDE  bool
}

func NewUPFunctionFeatures(bu, dd, dl, tr, ft, pf, he, treu, em, pd, ud, qu, trace, fr, pfde bool) *UPFunctionFeatures {

	return &UPFunctionFeatures{
		BUCP:  bu,
		DDND:  dd,
		DLBD:  dl,
		TRST:  tr,
		FTUP:  ft,
		PFDM:  pf,
		HEEU:  he,
		TREU:  treu,
		EMPU:  em,
		PDIU:  pd,
		UDBC:  ud,
		QUOAC: qu,
		TRACE: trace,
		FRRT:  fr,
		PFDE:  pfde,
	}

}
func (u *UPFunctionFeatures) Serialize() ([]byte, error) {
	var fByte, sByte uint8
	if u.BUCP {
		fByte = 1
	}
	if u.DDND {
		fByte |= uint8(2)
	}
	if u.DLBD {
		fByte |= uint8(4)
	}
	if u.TRST {
		fByte |= uint8(8)
	}
	if u.FTUP {
		fByte |= uint8(16)
	}
	if u.PFDM {
		fByte |= uint8(32)
	}
	if u.HEEU {
		fByte |= uint8(64)
	}
	if u.TREU {
		fByte |= uint8(128)
	}

	if u.EMPU {
		sByte = uint8(1)
	}

	if u.PDIU {
		sByte |= uint8(2)
	}
	if u.UDBC {
		sByte |= uint8(4)
	}
	if u.QUOAC {
		sByte |= uint8(8)
	}
	if u.TRACE {
		sByte |= uint8(16)
	}
	if u.FRRT {
		sByte |= uint8(32)
	}
	if u.PFDE {
		sByte |= uint8(64)
	}

	b := make([]byte, 2)
	b[0] = fByte
	b[1] = sByte
	return b, nil
}

func NewUPFunctionFeaturesFromByte(data []byte) *UPFunctionFeatures {
	if len(data) != 2 {
		return nil
	}

	bu := (uint8(data[0]&0x01) == 1)
	dd := (uint8(data[0]&0x02) == 2)
	dl := (uint8(data[0]&0x04) == 4)
	tr := (uint8(data[0]&0x08) == 8)
	ft := (uint8(data[0]&0x10) == 16)
	pf := (uint8(data[0]&0x20) == 32)
	he := (uint8(data[0]&0x40) == 64)
	treu := (uint8(data[0]&0x80) == 128)

	em := (uint8(data[1]&0x01) == 1)
	pd := (uint8(data[1]&0x02) == 2)
	ud := (uint8(data[1]&0x04) == 4)
	qu := (uint8(data[1]&0x08) == 8)
	trace := (uint8(data[1]&0x10) == 16)
	fr := (uint8(data[1]&0x20) == 32)
	pfde := (uint8(data[1]&0x40) == 64)

	return NewUPFunctionFeatures(bu, dd, dl, tr, ft, pf, he, treu, em, pd, ud, qu, trace, fr, pfde)

}
