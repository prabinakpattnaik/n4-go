package helper

import (
	"github.com/prabinakpattnaik/n4-go/ie"
)

type upipresourceinformation_setting struct {
	//TODO UP is here identified as a number. Interesting if key can be a single string or IP address
	setting map[int]*ie.UPIPResourceInformation
}
