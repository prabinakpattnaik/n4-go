package setting

import (
	"bitbucket.org/sothy5/n4-go/ie"
dt "github.com/fiorix/go-diameter/diam/datatype"
)

//Assign TunnelID gets input of UPIPResourceInformation, UE IP address, and assign new tunnel ID.
func Assign_tunnelID(u net.IP, tunnel_id uint32) (*ie.InformationElement,error) {
	//TODO manage tunnel_id against ue_ip_address *net.IPv4,
	// get IPv4 address from UPIPResourceInformation, tunnel_id and convet those info IE (FTEID)
	 fteid:=ie.NewFTEID(true,false, false, false,tunnel_id+1,u, nil,0)
	 b, err:=ftied.Serialize()
	 if err !=nil {
		 return nil, err
	 }
	 return &ie.NewInformationElement{
		 ie.IEFTEID,
		 0,
		 dt.OctetString(b)


	 }, nil
}
