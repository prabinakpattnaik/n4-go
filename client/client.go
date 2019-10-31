package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"bitbucket.org/sothy5/n4-go/msg"

	"bitbucket.org/sothy5/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

var (
	udpport           = 8805
	maxBufferSize     = 1024
	remoteIPv4address = net.IPv4(127, 0, 0, 1)
	sequanceNumber    = uint32(200)

	controlPlaneNodeID      = []byte{0x0, 0xC0, 0xa8, 0x1, 0x21}
	controlFunctionFeatures = []byte{0x00}
	PFCPMinHeaderSize       = 8
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
	pfcpHeader := msg.NewPFCPHeader(1, false, false, msg.AssociationSetupRequestType, length, 0, sequanceNumber, 0)

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

		pfcpAssociationSetupResponse, err := msg.FromPFCPMessage(pfcpMessage)
		if err != nil {
			fmt.Printf("error in FromPFCPMessage() %+v\n", err)
		}

		fmt.Printf("received message %+v\n", pfcpAssociationSetupResponse)
	}

}
