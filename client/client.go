package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"bitbucket.org/sothy5/n4-go/msg"

	setting "bitbucket.org/sothy5/n4-go/client/internal/helper"
	"bitbucket.org/sothy5/n4-go/client/internal/session"
	"bitbucket.org/sothy5/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
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
		log.Printf(" resolveUDP Addr err %+v\n", err)
	}

	conn, err := net.DialUDP("udp4", nil, dst)
	if err != nil {
		//TODO handle this error
		log.Printf("failure in connection setup %+v\n", err)
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
		fmt.Println("Nil conn pointer")
	}
	_, err := c.Conn.Write(b)
	if err != nil {
		log.Printf("failure in writing over connection: %+v\n", err)
		return err
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
			log.Print(err)
		}
		fmt.Printf("received pfcpSessionEstablishmentResponse [%x]\n", rb)
		pfcpMessage, err := msg.MessageFromBytes(rb)
		if err != nil {
			log.Print(err)
		}
		pfcp, err := msg.FromPFCPMessage(pfcpMessage)
		if err != nil {
			fmt.Printf("error in FromPFCPMessage() %+v\n", err)
		}
		pfcpSessionEstablishmentResponse, ok := pfcp.(msg.PFCPSessionEstablishmentResponse)
		if ok {
			fmt.Printf("received pfcpSessionEstablishmentResponse %+v\n", pfcpSessionEstablishmentResponse)
			fmt.Printf("received pfcpSessionEstablishmentResponse Header%+v\n", pfcpSessionEstablishmentResponse.Header)

			sessionRequestResponse := session.SessionRequestResponse{
				SResponse: &pfcpSessionEstablishmentResponse,
			}

			sessionEntity.Inc(pfcpSessionEstablishmentResponse.Header.SequenceNumber, sessionRequestResponse)
			fmt.Printf("sequence number is equal for request and response %d\n", pfcpSessionEstablishmentResponse.Header.SequenceNumber)
		} else {
			fmt.Printf("Error")
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
			fmt.Printf("error in MessageFromBytes() %+v\n", err)
		}

		pfcp, err := msg.FromPFCPMessage(pfcpMessage)
		if err != nil {
			fmt.Printf("error in FromPFCPMessage() %+v\n", err)
		}
		pfcpAssociationSetupResponse, ok := pfcp.(msg.PFCPAssociationSetupResponse)
		if !ok {
			fmt.Printf("wrong in  type assertation")
		}
		fmt.Printf("received message for UserPlaneIPResourceInformation %+v\n", pfcpAssociationSetupResponse.UserPlaneIPResourceInformation)
		b, err := pfcpAssociationSetupResponse.UserPlaneIPResourceInformation.Serialize()
		if err != nil {
			fmt.Printf("error in pfcpAssociationSetupResponse.UserPlaneIPResourceInformation.Serialize() %+v\n", err)
		}
		UPIPResourceInformation := ie.NewUPIPResourceInformationFromByte(pfcpAssociationSetupResponse.UserPlaneIPResourceInformation.Length, b[4:])
		fmt.Printf("received UPIPResourceInformation %+v\n", UPIPResourceInformation)

		//setting := make(map[int]*ie.UPIPResourceInformation)
		//setting[1] = UPIPResourceInformation

		go RecvProcess(client)

		for i := 0; i < 10; i++ {
			teid++
			sequenceNumber++
			seid++
			time.Sleep(2 * time.Second)

			fteid, err := setting.Assign_tunnelID(UPIPResourceInformation.IPv4Address, teid)
			pfcpSessionEstablishmentRequest, err := session.CreateNewSession(seid, sequenceNumber, nodeIP, seid, 1, 1, 0, fteid, 2, 1)
			if err != nil {
				fmt.Printf("error in pfcpSessionEstablishmentRequest %+v\n", err)
				continue

			}

			b, err := pfcpSessionEstablishmentRequest.Serialize()
			if err != nil {
				fmt.Printf("error in pfcpSessionEstablishmentRequest %+v\n", err)
				continue
			}
			sessionRequestResponse := session.SessionRequestResponse{
				SRequest: pfcpSessionEstablishmentRequest,
			}
			sessionEntity.Inc(sequenceNumber, sessionRequestResponse)
			client.Write(b)

		}
		time.Sleep(2 * time.Second)

		//TODO: Keep NodeID, UPFunctionFeatures, and UPIPResourceInformation

	}

}
