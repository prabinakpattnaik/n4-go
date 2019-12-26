package main

import (
	"log"
	"net"

	"bitbucket.org/sothy5/n4-go/msg"
	"bitbucket.org/sothy5/n4-go/server/internal/session"
	"bitbucket.org/sothy5/n4-go/util/server"
)

var (
	udpport   = 8805
	ipaddress = "127.0.0.1"
	seid      = uint64(10000)
	nodeIP    = net.ParseIP("127.0.0.1")

	seidEntity = session.SEIDEntity{M: make(map[uint64]uint64)}
)

func handler(conn net.PacketConn, peer net.Addr, m *msg.PFCPMessage) {
	// this function will just print the received PFCP message, without replying

	switch m.Header.MessageType {
	case msg.HeartbeatRequestType:
		log.Print("Not handled Heartbeat Request type")
	case msg.AssociationSetupRequestType:
		b, err := msg.ProcessAssociationSetupRequest(m)
		if err == nil {
			if _, err := conn.WriteTo(b, peer); err != nil {
				log.Printf("Cannot send Message Type {%d} to client: %v", msg.AssociationSetupResponseType, err)
			}
		}
	case msg.SessionEstablishmentRequestType:
		b, err := msg.ProcessPFCPSessionEstablishmentRequest(m, nodeIP, seid)
		if err == nil {
			if _, err := conn.WriteTo(b, peer); err != nil {
				log.Printf("Cannot send Message Type {%d} to client: %v", msg.SessionEstablishmentResponseType, err)
			}
		}
		seidEntity.Inc(m.Header.SessionEndpointIdentifier, seid)
		seid = seid + 1

	case msg.SessionModificationRequestType:
		sSEID := m.Header.SessionEndpointIdentifier
		if sSEID > 0 {
			b, err := msg.ProcessPFCPSessionModificationRequest(m, sSEID)
			if err == nil {
				if _, err := conn.WriteTo(b, peer); err != nil {
					log.Printf("Cannot send Message Type {%d} to client: %v", msg.SessionModificationResponseType, err)
				}
			}

		} else {
			log.Printf("Not valid SEID in PFCPMessage: %+v\n", m)
		}
	case msg.SessionDeletionRequestType:
		sSEID := m.Header.SessionEndpointIdentifier
		if sSEID > 0 {

			b, err := msg.ProcessPFCPSessionDeletionRequest(m, sSEID)
			if err == nil {
				if _, err := conn.WriteTo(b, peer); err != nil {
					log.Printf("Cannot send Message Type {%d} to client: %v", msg.SessionModificationResponseType, err)
				}
			}

		} else {
			log.Printf("Not valid SEID in PFCPMessage: %+v\n", m)
		}

	default:
		log.Printf("Not Handled this PFCPMessage: %+v\n", m)

	}

}

func main() {
	laddr := net.UDPAddr{
		IP:   net.ParseIP(ipaddress),
		Port: udpport,
	}
	server := server.NewServer(laddr, handler)
	defer server.Close()
	if err := server.ActivateAndServe(); err != nil {
		log.Panic(err)
	}
}
