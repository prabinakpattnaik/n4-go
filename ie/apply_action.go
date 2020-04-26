package ie

type ApplyAction struct {
	Dupl bool
	Nocp bool
	Buff bool
	Forw bool
	Drop bool
}

func (a *ApplyAction) Serialize() ([]byte, error) {
	var b uint8
	if a.Drop {
		b = 1
	}
	if a.Forw {
		b = 2
	}
	if a.Buff {
		b = 4
	}
	if a.Buff && a.Nocp {
		b = 24
	}
	if a.Dupl {
		b |= 32
	}

	return []byte{b}, nil
}
