package server_wrap

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/prabinakpattnaik/n4-go/ie"
	"github.com/prabinakpattnaik/n4-go/msg"
	"github.com/prabinakpattnaik/n4-go/util/se"
	"github.com/prabinakpattnaik/n4-go/util/server"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

var (
	controlFunctionFeatures = []byte{0x00}
	controlPlaneNodeID      []byte
	cha                     chan []ie.UPIPResourceInformation
	ftupCha                 chan bool
	cpSEIDDPSEID            = &se.CPSEIDDPSEIDEntity{M: make(map[uint64]uint64)}
)

func handler(conn net.PacketConn, peer net.Addr, m *msg.PFCPMessage) {
	// this function will just print the received PFCP message, without replying

	switch m.Header.MessageType {
	case msg.HeartbeatRequestType:
		pfcp, err := msg.FromPFCPMessage(m)

		if err != nil {
			log.Printf("error in casting %v", err)
			return
		}
		pfcpHeartBeat, ok := pfcp.(*msg.Heartbeat)
		if ok {
			r := ie.NewInformationElement(
				ie.IERecoveryTimestamp,
				0,
				dt.Time(time.Now()),
			)
			h := pfcpHeartBeat.GetHeader()
			h.MessageType = msg.HeartbeatResponseType
			h.MessageLength = msg.PFCPBasicMessageSize + ie.IEBasicHeaderSize + r.Len()
			heartbeat := msg.NewHeartbeat(h, &r)
			b, _ := heartbeat.Serialize()
			conn.WriteTo(b, peer)
		}

	case msg.AssociationSetupRequestType:
		b, err := ProcessAssociationSetupRequest(m)
		if err == nil {
			if _, err := conn.WriteTo(b, peer); err != nil {
				log.Printf("Cannot send Message Type {%d} to client: %v", msg.AssociationSetupResponseType, err)
			}
		}
	case msg.SessionReportRequestType:

	case msg.SessionEstablishmentRequestType:
		log.Print("Not handled  SessionEstablishmentRequestType")

	case msg.SessionModificationRequestType:
		log.Print("Not handled  SessionModificationRequestType")

	case msg.SessionDeletionRequestType:
		log.Print("Not handled  SessionDeletionRequestType")

	default:
		log.Printf("Not Handled this PFCPMessage: %+v\n", m)

	}

}

func ProcessAssociationSetupRequest(m *msg.PFCPMessage) ([]byte, error) {
	pfcp, err := msg.FromPFCPMessage(m)
	if err != nil {
		return nil, err
	}
	pfcpAssociationSetupRequest, ok := pfcp.(msg.PFCPAssociationSetupRequest)
	if !ok {
		log.Print("Not received pfcpAssociationSetuprequest")
		return nil, fmt.Errorf("Not received PFCPAssocationSetup Request")
	}
	if pfcpAssociationSetupRequest.NodeID == nil {
		//TODO not handled propelly
		return nil, fmt.Errorf("No valid NodeID ")
	}
	if pfcpAssociationSetupRequest.RecoveryTimeStamp == nil {
		return nil, fmt.Errorf("No RecoveryTimestamp")
	}

	var upIPRIs []ie.UPIPResourceInformation
	for _, informationElement := range pfcpAssociationSetupRequest.UserPlaneIPResourceInformation {
		b, _ := informationElement.Serialize()
		upIPResourceInformation := ie.NewUPIPResourceInformationFromByte(informationElement.Len(), b[4:])
		upIPRIs = append(upIPRIs, *upIPResourceInformation)
	}
	if len(upIPRIs) > 0 {
		cha <- upIPRIs
	}
	informationElement := pfcpAssociationSetupRequest.UPFunctionFeatures
	b, _ := informationElement.Serialize()

	upFunctionFeatures := ie.NewUPFunctionFeaturesFromByte(b[4:])
	if upFunctionFeatures.FTUP {
		if len(upIPRIs) == 0 {
			ftupCha <- upFunctionFeatures.FTUP
		}
	}

	n := ie.NewInformationElement(
		ie.IENodeID,
		0,
		dt.OctetString(controlPlaneNodeID),
	)
	length := n.Len() + ie.IEBasicHeaderSize

	c := ie.NewInformationElement(
		ie.IECause,
		0,
		dt.OctetString([]byte{0x01}),
	)
	length = length + c.Len() + ie.IEBasicHeaderSize

	r := ie.NewInformationElement(
		ie.IERecoveryTimestamp,
		0,
		dt.Time(time.Now()),
	)
	length = length + r.Len() + ie.IEBasicHeaderSize

	cp := ie.NewInformationElement(
		ie.IECPFunctionFeatures,
		0,
		dt.OctetString(controlFunctionFeatures),
	)
	length = length + cp.Len() + ie.IEBasicHeaderSize

	length = length + msg.PFCPBasicMessageSize
	pfcpHeader := msg.NewPFCPHeader(1, false, false, msg.AssociationSetupResponseType, length, 0, pfcpAssociationSetupRequest.GetHeader().SequenceNumber, 0)

	ar := msg.NewPFCPAssociationSetupResponse(pfcpHeader, &n, &c, &r, nil, &cp, nil)
	return ar.Serialize()
}

func Run(ipaddress string, udpport int, cpNID []byte, ch chan []ie.UPIPResourceInformation, ftupC chan bool, cPSEIDDPSEID *se.CPSEIDDPSEIDEntity) {
	laddr := net.UDPAddr{
		IP:   net.ParseIP(ipaddress),
		Port: udpport,
	}
	cpSEIDDPSEID = cPSEIDDPSEID
	controlPlaneNodeID = cpNID
	cha = ch
	ftupCha = ftupC
	server := server.NewServer(laddr, handler)
	defer server.Close()
	if err := server.ActivateAndServe(); err != nil {
		log.Panic(err)
	}
}
