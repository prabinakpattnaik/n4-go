package util_3gpp

type Dnn []uint8

func (d *Dnn) Serialize() (data []byte, err error) {

	data = append(data, uint8(len(*d)))
	data = append(data, (*d)...)

	return data, nil
}

func (d *Dnn) DnnFromBinary(data []byte) error {

	(*d) = data[1:]
	return nil
}
