package main

import (
	"fmt"
	"net"
	"time"

	"bitbucket.org/sothy5/n4-go/msg"

	setting "bitbucket.org/sothy5/n4-go/client/internal/helper"
	"bitbucket.org/sothy5/n4-go/client/internal/session"
	"bitbucket.org/sothy5/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	log "github.com/sirupsen/logrus"
)

var (
	udpport           = 8805
	maxBufferSize     = 1024
	remoteIPv4address = net.IPv4(127, 0, 0, 1)
	sequenceNumber    = uint32(100)
	seid              = uint64(100)

	controlPlaneNodeID         = []byte{0x0, 0xC0, 0xa8, 0x1, 0x20}
	nodeIP                     = net.ParseIP("192.168.1.32")
	controlFunctionFeatures    = []byte{0x00}
	PFCPMinHeaderSize          = 8
	UPIPResourceInformationMap map[int]*ie.UPIPResourceInformation
	teid                       uint32 = 0
	sessionEntity                     = session.SessionEntity{M: make(map[uint32]session.SessionRequestResponse)}
	seidsnEntity                      = session.SEIDSNEntity{M: make(map[uint64]session.SNCollection)}
)

// Client implements a PFCP client
type Client struct {
	LocalAddr  net.Addr
	RemoteAddr *net.UDPAddr
	Conn       *net.UDPConn
}

// NewClient returns a Client with default settings
func NewClient() *Client {
	raddr := fmt.Sprintf("%s:%d", remoteIPv4address, udpport)

	dst, err := net.ResolveUDPAddr("udp", raddr)

	if err != nil {
		log.WithError(err).Error("resolveUDP Addr err ")
	}

	conn, err := net.DialUDP("udp4", nil, dst)
	if err != nil {
		//TODO handle this error
		log.WithError(err).Error("failure in connection setup")
		return nil
	}
	return &Client{
		RemoteAddr: dst,
		Conn:       conn,
	}
}

func (c *Client) Read() ([]byte, error) {
	buffer := make([]byte, maxBufferSize)

	nRead, _, err := c.Conn.ReadFrom(buffer)
	if err != nil {
		return nil, err
	}
	b := make([]byte, nRead)
	copy(b, buffer)
	return b, nil

}

func (c *Client) Write(b []byte) error {
	if c.Conn == nil {
		log.Fatal("Nil Conn pointer")
	}
	_, err := c.Conn.Write(b)
	if err != nil {
		log.WithError(err).Fatal("Not possible to write over Conn")

	}
	return nil

}

//Close the UDP connection
func (c *Client) Close() {
	c.Conn.Close()

}

func RecvProcess(c *Client) {
	for {
		rb, err := c.Read()
		if err != nil {
			log.WithError(err).Fatal("connection reading error")
		}

		pfcpMessage, err := msg.MessageFromBytes(rb)
		if err != nil {
			log.WithError(err).Info("Error in received pfcpSessionEstablishmentResponse")
		}
		pfcp, err := msg.FromPFCPMessage(pfcpMessage)
		if err != nil {
			log.WithError(err).Info("error in FromPFCPMessage")
		}

		pfcpSessionEstablishmentResponse, ok := pfcp.(msg.PFCPSessionEstablishmentResponse)
		if ok {
			log.WithFields(log.Fields{"data": rb}).Info("received pfcpSessionEstablishmentResponse")
			sessionRequestResponse := session.SessionRequestResponse{
				SResponse: &pfcpSessionEstablishmentResponse,
			}
			sessionEntity.Inc(pfcpSessionEstablishmentResponse.Header.SequenceNumber, sessionRequestResponse)

		}
		//pfcpSessionModificationResponse
		pfcpSessionModificationResponse, ok := pfcp.(msg.PFCPSessionModificationResponse)
		if ok {
			log.WithFields(log.Fields{"data": rb}).Info("received pfcpSessionModificationResponse")
			sessionRequestResponse := session.SessionRequestResponse{
				SResponse: &pfcpSessionModificationResponse,
			}
			sessionEntity.Inc(pfcpSessionModificationResponse.Header.SequenceNumber, sessionRequestResponse)

		}

		//pfcpSessionModificationResponse
		pfcpSessionDeletionResponse, ok := pfcp.(msg.PFCPSessionDeletionResponse)
		if ok {
			log.WithFields(log.Fields{"data": rb}).Info("received pfcpSessionDeletionResponse")
			sessionRequestResponse := session.SessionRequestResponse{
				SResponse: &pfcpSessionDeletionResponse,
			}
			sessionEntity.Inc(pfcpSessionDeletionResponse.Header.SequenceNumber, sessionRequestResponse)

		}

	}

}

func main() {

	client := NewClient()

	//TODO: create HeartBeat Request
	//Create PFCPAssociationSetupRequest
	//How to correlate  the reply with request

	//How to increase the sequence number

	//Perform PFCPAssociationSetupRequest only

	var length uint16

	n := ie.NewInformationElement(
		ie.IENodeID,
		0,
		dt.OctetString(controlPlaneNodeID),
	)
	length = n.Len() + ie.IEBasicHeaderSize

	r := ie.NewInformationElement(
		ie.IERecoveryTimestamp,
		0,
		dt.Time(time.Now()),
	)
	length = length + r.Len() + ie.IEBasicHeaderSize

	c := ie.NewInformationElement(
		ie.IECPFunctionFeatures,
		0,
		dt.OctetString(controlFunctionFeatures),
	)
	length = length + c.Len() + ie.IEBasicHeaderSize

	length = length + msg.PFCPBasicMessageSize
	pfcpHeader := msg.NewPFCPHeader(1, false, false, msg.AssociationSetupRequestType, length, 0, sequenceNumber, 0)

	ar := msg.NewPFCPAssociationSetupRequest(pfcpHeader, &n, &r, nil, &c, nil)
	b, _ := ar.Serialize()
	if len(b) > PFCPMinHeaderSize {
		client.Write(b)
		rb, err := client.Read()
		if err != nil {
			log.Print(err)
		}
		pfcpMessage, err := msg.MessageFromBytes(rb)

		if err != nil {
			log.WithError(err).Fatal("error from MessageFromBytes")

		}

		pfcp, err := msg.FromPFCPMessage(pfcpMessage)
		if err != nil {
			log.WithError(err).Fatal("error from FromPFCPMessage")
		}
		pfcpAssociationSetupResponse, ok := pfcp.(msg.PFCPAssociationSetupResponse)
		if !ok {
			log.WithError(err).Fatal("wrong type asseration of PFCPAssociationSetupResponse")
		}
		log.WithFields(log.Fields{"UserPlaneIPResourceInformation": pfcpAssociationSetupResponse.UserPlaneIPResourceInformation}).Info("Received Information")
		b, err := pfcpAssociationSetupResponse.UserPlaneIPResourceInformation.Serialize()
		if err != nil {
			log.WithError(err).Fatal("error in pfcpAssociationSetupResponse.UserPlaneIPResourceInformation.Serialize")

		}
		UPIPResourceInformation := ie.NewUPIPResourceInformationFromByte(pfcpAssociationSetupResponse.UserPlaneIPResourceInformation.Length, b[4:])
		log.WithFields(log.Fields{"V4": UPIPResourceInformation.V4,
			"V6":              UPIPResourceInformation.V6,
			"TEIDRI":          UPIPResourceInformation.TEIDRI,
			"ASSONI":          UPIPResourceInformation.ASSONI,
			"ASSOSI":          UPIPResourceInformation.ASSOSI,
			"TEIDRange":       UPIPResourceInformation.TEIDRange,
			"IPv4Address":     UPIPResourceInformation.IPv4Address,
			"IPv6Address":     UPIPResourceInformation.IPv6Address,
			"NetworkInstance": UPIPResourceInformation.NetworkInstance,
			"SourceInterface": UPIPResourceInformation.SourceInterface,
		}).Info("Received UserPlaneIPResourceInformation")

		//setting := make(map[int]*ie.UPIPResourceInformation)
		//setting[1] = UPIPResourceInformation

		go RecvProcess(client)

		for i := 0; i < 10; i++ {
			teid++
			sequenceNumber++
			seid++
			time.Sleep(2 * time.Second)

			fteid, err := setting.Assign_tunnelID(UPIPResourceInformation.IPv4Address, teid)
			pfcpSessionEstablishmentRequest, err := session.CreateSession(seid, sequenceNumber, nodeIP, seid, 1, 1, 0, fteid, 2, 1)
			if err != nil {
				log.WithError(err).Error("error in pfcpSessionEstablishmentRequest")
				continue

			}

			b, err := pfcpSessionEstablishmentRequest.Serialize()
			if err != nil {
				log.WithError(err).Error("error in pfcpSessionEstablishmentRequest serialization")
				continue
			}
			sessionRequestResponse := session.SessionRequestResponse{
				SRequest: pfcpSessionEstablishmentRequest,
			}
			sessionEntity.Inc(sequenceNumber, sessionRequestResponse)
			seidsnEntity.Inc(seid, 1, sequenceNumber)
			client.Write(b)

		}
		time.Sleep(5 * time.Second)

		//TODO: Keep NodeID, UPFunctionFeatures, and UPIPResourceInformation
		//if we knew SEID,
		//PDR ID, FAR ID ? unique within a session or anytime?
		ueIPAddress := net.ParseIP("10.1.1.1")
		rIPAddress := net.ParseIP("192.168.1.1")
		rteid := uint32(5000)
		sn := uint32(101)
		for i := 0; i < 10; i++ {
			srr := sessionEntity.Value(sn)
			sequenceNumber++
			if srr.SRequest != nil {
				smr, err := session.ModifySession(srr.SRequest.GetHeader().SessionEndpointIdentifier, sequenceNumber, 2, 2, ie.Core, ueIPAddress, rteid, rIPAddress, uint8(ie.FORW), ie.Access)
				if err != nil {
					log.WithError(err).Error("error in pfcpSessionModificationRequest")
					continue
				}
				b, err := smr.Serialize()
				if err != nil {
					log.WithError(err).Error("error in pfcpSessionModificationRequest serialization")
					continue
				}
				client.Write(b)
				sessionRequestResponse := session.SessionRequestResponse{
					SRequest: smr,
				}
				sessionEntity.Inc(sequenceNumber, sessionRequestResponse)
				seidsnEntity.Inc(srr.SRequest.GetHeader().SessionEndpointIdentifier, 2, sequenceNumber)
				sn++
			}
		}
		time.Sleep(7 * time.Second)

		sn = uint32(101)
		for i := 0; i < 10; i++ {
			srr := sessionEntity.Value(sn)
			sequenceNumber++
			if srr.SRequest != nil {
				pfcpHeader := msg.NewPFCPHeader(1, false, true, msg.SessionDeletionRequestType, 12, srr.SRequest.GetHeader().SessionEndpointIdentifier, sequenceNumber, 0)
				b := pfcpHeader.Serialize()
				client.Write(b)
				log.WithFields(log.Fields{"data": b}).Info("received pfcpSessionDeletionRequest")

				sn++

				sdr := msg.NewPFCPSessionDeletionRequest(pfcpHeader)
				sessionRequestResponse := session.SessionRequestResponse{
					SRequest: &sdr,
				}
				sessionEntity.Inc(sequenceNumber, sessionRequestResponse)
				seidsnEntity.Inc(srr.SRequest.GetHeader().SessionEndpointIdentifier, 3, sequenceNumber)
			}
		}
		time.Sleep(5 * time.Second)

	}

}
