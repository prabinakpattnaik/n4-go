package flowdesc

type IPFilterRule struct {
	action   bool   // true: permit, false: deny
	dir      bool   // false: in, true: out
	proto    int    // protocal number
	srcIp    string // <address/mask>
	srcPorts string // [ports]
	dstIp    string // <address/mask>
	dstPorts string // [ports]
}

func (i *IPFilterRule) Serilize() ([]byte, error) {
	return nil, nil
}

//Flow description.
// flow_desc="permit out ip from any to assigned"
//flow_desc="permit out ip from any to assigned"
