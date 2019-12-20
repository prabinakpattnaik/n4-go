package sr

import (
	"bitbucket.org/sothy5/n4-go/ie"
)

// EthernetTrafficInformation
type EthernetTrafficInformation struct {
	MACAddressesDetected *ie.InformationElement
	MADAddressesRemoved  *ie.InformationElement
}
