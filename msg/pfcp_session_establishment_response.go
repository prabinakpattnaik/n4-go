package msg

import (
	"bitbucket.org/sothy5/n4-go/ie"
)

//PFCPSessionEstablishmentResponse
type PFCPSessionEstablishmentResponse struct {
	Header                     *PFCPHeader
	NodeID                     *ie.InformationElement
	Cause                      *ie.InformationElement
	OffendingIE                *ie.InformationElement
	UPFSEID                    *ie.InformationElement
	LoadControlInformation     *ie.InformationElement
	OverloadControlInformation *ie.InformationElement
	FailedRuleID               *ie.InformationElement
	CreaatedTrafficEndpoint    *ie.InformationElement
}
