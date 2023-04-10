package helper

import (
	"net"

	"github.com/prabinakpattnaik/n4-go/ie"
)

//Assign TunnelID gets input of UPIPResourceInformation, UE IP address, and assign new tunnel ID.
func Assign_tunnelID(u net.IP, tunnel_id uint32) (*ie.FTEID, error) {
	//TODO manage tunnel_id against ue_ip_address *net.IPv4,
	// get IPv4 address from UPIPResourceInformation, tunnel_id and convet those info IE (FTEID)
	var fteid *ie.FTEID
	if u == nil && tunnel_id == 0 {
		fteid = ie.NewFTEID(true, false, true, false, 0, nil, nil, 0)
	} else {
		fteid = ie.NewFTEID(true, false, false, false, tunnel_id, u, nil, 0)
	}
	return fteid, nil

}
