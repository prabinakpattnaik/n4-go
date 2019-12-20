package sr

import (
	"bitbucket.org/sothy5/n4-go/ie"
)

// ApplicationDetectionInformation
type ApplicationDetectionInformation struct {
	ApplicationID         *ie.InformationElement
	ApplicationInstanceID *ie.InformationElement
	FlowInformation       *ie.InformationElement
}
