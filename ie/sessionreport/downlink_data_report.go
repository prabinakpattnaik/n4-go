package sr

import "bitbucket.org/sothy5/n4-go/ie"

type DownlinkDataReport struct {
	PDRID                          ie.InformationElements
	DownlinkDataServiceInformation ie.InformationElements
}
