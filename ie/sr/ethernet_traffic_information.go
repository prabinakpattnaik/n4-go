package sr

import (
	"github.com/prabinakpattnaik/n4-go/ie"
)

// EthernetTrafficInformation
type EthernetTrafficInformation struct {
	MACAddressesDetected *ie.InformationElement
	MADAddressesRemoved  *ie.InformationElement
}
