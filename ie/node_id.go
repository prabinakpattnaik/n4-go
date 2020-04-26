package ie
import (
"net"
)

const (
	NodeIdTypeIpv4Address uint8 = iota
	NodeIdTypeIpv6Address
	NodeIdTypeFqdn
)

type NodeID struct {
	NodeIdType  uint8 // 0x00001111
	NodeIdValue []byte
}

func NewNodeID(nodeIdType uint8, nodeIdValue []byte) *NodeID {
	return &NodeID{
		NodeIdType:  nodeIdType,
		NodeIdValue: nodeIdValue,
	}
}

func NewNodeIDFromBytes(length uint16, b []byte) *NodeID {
	if length == 0 {
		return nil
	}

	firstByte := b[0]
	nodeIDType := uint8(firstByte)
	var nodeIdValue []byte
	if nodeIDType == 0 && length == 5 {
		nodeIdValue = make([]byte, 4)
		copy(nodeIdValue, b[1:5])
	} else if nodeIDType == 1 {
		return nil
	} else if nodeIDType == 2 {
		return nil
	}

	return &NodeID{
		NodeIdType:  nodeIDType,
		NodeIdValue: nodeIdValue,
	}

}

func (n *NodeID) ResolveNodeIdToIp() net.IP {
     if n.NodeIdType == NodeIdTypeIpv4Address {
	     return n.NodeIdValue
     }else if n.NodeIdType == NodeIdTypeIpv6Address {
	     return n.NodeIdValue
     
     }
     return nil
}

func (n *NodeID) Serialize() ([]byte, error) {
	b := make([]byte, 1)
	b[0] = n.NodeIdType
	b = append(b, n.NodeIdValue...)
	return b, nil
}
