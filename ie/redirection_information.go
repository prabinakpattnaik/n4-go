package ie

type RedirectInformation struct {
	RedirectAddressType         uint8 // 0x00001111
	RedirectServerAddressLength uint16
	RedirectServerAddress       []byte
}
