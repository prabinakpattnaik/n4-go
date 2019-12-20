package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"bitbucket.org/sothy5/n4-go/ie/urr"
	"bitbucket.org/sothy5/n4-go/msg"

	setting "bitbucket.org/sothy5/n4-go/client/internal/helper"
	"bitbucket.org/sothy5/n4-go/client/internal/server_wrap"
	"bitbucket.org/sothy5/n4-go/client/internal/session"
	"bitbucket.org/sothy5/n4-go/client/internal/usage_report"
	"bitbucket.org/sothy5/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	udpport       = 8805
	maxBufferSize = 1024

	sequenceNumber = uint32(100)
	seid           = uint64(100)

	controlFunctionFeatures    = []byte{0x00}
	PFCPMinHeaderSize          = 8
	UPIPResourceInformationMap map[int]*ie.UPIPResourceInformation
	teid                       uint32 = 0
	sessionEntity                     = session.SessionEntity{M: make(map[uint32]session.SessionRequestResponse)}
	seidsnEntity                      = session.SEIDSNEntity{M: make(map[uint64]session.SNCollection)}
	NetworkInstance                   = "epc"
	seSEIDSN                          = session.SESEIDSNEntity{M: make(map[uint64]uint32)}
	cpSEIDDPSEID                      = session.CPSEIDDPSEIDEntity{M: make(map[uint64]uint64)}
)

// Client implements a PFCP client
type Client struct {
	LocalAddr  net.Addr
	RemoteAddr *net.UDPAddr
	Conn       *net.UDPConn
}

// NewClient returns a Client with default settings
func NewClient(rAddress string, lAddress string) *Client {
	remoteIPv4address := net.ParseIP(rAddress)
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
		if len(rb) > 0 {
			pfcpMessage, err := msg.MessageFromBytes(rb)
			if err != nil {
				log.WithError(err).Info("Error in received pfcpMessage")
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
				// success verification
				if pfcpSessionEstablishmentResponse.UPFSEID != nil {
					b, _ := pfcpSessionEstablishmentResponse.UPFSEID.Serialize()
					fseid := ie.NewFSEIDFromByte(b[4:])
					log.WithFields(log.Fields{"fseid v4": fseid.V4,
						"fseid v4 address": fseid.IP4Address,
						"fseid seid":       fseid.SEID,
					}).Info("received DP FSEID")
					//SessionEstablishment Request and Response have SEID
					cpSEIDDPSEID.Inc(pfcpSessionEstablishmentResponse.GetHeader().SessionEndpointIdentifier, fseid.SEID)

				}

				//SEID from client, new SEID to be used by client, upon assigned by DP.
				continue

			}
			//pfcpSessionModificationResponse
			pfcpSessionModificationResponse, ok := pfcp.(msg.PFCPSessionModificationResponse)
			if ok {
				log.WithFields(log.Fields{"data": rb}).Info("received pfcpSessionModificationResponse")
				sessionRequestResponse := session.SessionRequestResponse{
					SResponse: &pfcpSessionModificationResponse,
				}
				sessionEntity.Inc(pfcpSessionModificationResponse.Header.SequenceNumber, sessionRequestResponse)
				continue
			}

			//pfcpSessionModificationResponse
			pfcpSessionDeletionResponse, ok := pfcp.(msg.PFCPSessionDeletionResponse)
			if ok {
				log.WithFields(log.Fields{"data": rb}).Info("received pfcpSessionDeletionResponse")
				sessionRequestResponse := session.SessionRequestResponse{
					SResponse: &pfcpSessionDeletionResponse,
				}
				sessionEntity.Inc(pfcpSessionDeletionResponse.Header.SequenceNumber, sessionRequestResponse)
				continue
			}

			pfcpHeartbeat, ok := pfcp.(*msg.Heartbeat)
			if ok {
				log.WithFields(log.Fields{"data": rb}).Info("received PFCPHeartbeat Request")
				r := ie.NewInformationElement(
					ie.IERecoveryTimestamp,
					0,
					dt.Time(time.Now()),
				)
				h := pfcpHeartbeat.GetHeader()
				h.MessageType = msg.HeartbeatResponseType
				h.MessageLength = msg.PFCPBasicMessageSize + ie.IEBasicHeaderSize + r.Len()
				heartbeat := msg.NewHeartbeat(h, &r)
				b, _ := heartbeat.Serialize()
				c.Write(b)
				continue
			}

			log.WithFields(log.Fields{"data": rb}).Info("something went")
		}
	}

}

func run(c *cli.Context) error {
	lIPv4address := c.String("localIP")
	nodeIP := net.ParseIP(lIPv4address)
	controlPlaneNodeID := []byte{0x0}
	controlPlaneNodeID = append(controlPlaneNodeID, nodeIP.To4()...)

	rIPv4address := c.String("remoteIP")
	client := NewClient(rIPv4address, lIPv4address)
	var upIPRI ie.UPIPResourceInformation
	var upIPRIs []ie.UPIPResourceInformation
	upIPRIC := make(chan []ie.UPIPResourceInformation)
	ftupC := make(chan bool)
	var ftup bool

	if !c.Bool("clientInit") {
		go server_wrap.Run(lIPv4address, udpport, controlPlaneNodeID, upIPRIC, ftupC)

		select {
		case upIPRIs = <-upIPRIC:

		case ftup = <-ftupC:
			break
		}

		time.Sleep(20 * time.Second)

	} else {

		time.Sleep(2 * time.Second)
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

		cp := ie.NewInformationElement(
			ie.IECPFunctionFeatures,
			0,
			dt.OctetString(controlFunctionFeatures),
		)
		length = length + cp.Len() + ie.IEBasicHeaderSize

		length = length + msg.PFCPBasicMessageSize
		pfcpHeader := msg.NewPFCPHeader(1, false, false, msg.AssociationSetupRequestType, length, 0, sequenceNumber, 0)

		ar := msg.NewPFCPAssociationSetupRequest(pfcpHeader, &n, &r, nil, &cp, nil)
		b, _ := ar.Serialize()
		if len(b) > PFCPMinHeaderSize {
			client.Write(b)
			//TODO, it is not blocking call.
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

			for _, informationElement := range pfcpAssociationSetupResponse.UserPlaneIPResourceInformation {
				b, _ := informationElement.Serialize()
				upIPResourceInformation := ie.NewUPIPResourceInformationFromByte(informationElement.Len(), b[4:])
				log.WithFields(log.Fields{"V4": upIPResourceInformation.V4,
					"V6":              upIPResourceInformation.V6,
					"TEIDRI":          upIPResourceInformation.TEIDRI,
					"ASSONI":          upIPResourceInformation.ASSONI,
					"ASSOSI":          upIPResourceInformation.ASSOSI,
					"TEIDRange":       upIPResourceInformation.TEIDRange,
					"IPv4Address":     upIPResourceInformation.IPv4Address,
					"IPv6Address":     upIPResourceInformation.IPv6Address,
					"NetworkInstance": upIPResourceInformation.NetworkInstance,
					"SourceInterface": upIPResourceInformation.SourceInterface,
				}).Info("Received UserPlaneIPResourceInformation")
				upIPRIs = append(upIPRIs, *upIPResourceInformation)
			}
			//setting := make(map[int]*ie.UPIPResourceInformation)
			//setting[1] = UPIPResourceInformation

		}
	}

	l := len(upIPRIs)
	if l == 1 {
		upIPRI = upIPRIs[0]
	} else if l > 1 {
		for _, u := range upIPRIs {
			if string(u.NetworkInstance[1:]) == NetworkInstance {
				upIPRI = u

			}
		}
	}

	if l > 0 {
		log.WithFields(log.Fields{"V4": upIPRI.V4,
			"V6":              upIPRI.V6,
			"TEIDRI":          upIPRI.TEIDRI,
			"ASSONI":          upIPRI.ASSONI,
			"ASSOSI":          upIPRI.ASSOSI,
			"TEIDRange":       upIPRI.TEIDRange,
			"IPv4Address":     upIPRI.IPv4Address,
			"IPv6Address":     upIPRI.IPv6Address,
			"NetworkInstance": upIPRI.NetworkInstance,
			"SourceInterface": upIPRI.SourceInterface,
		}).Info("selected  UserPlaneIPResourceInformation for TEID allocation")
		if upIPRI.TEIDRI > 0 {
			//TODO
			v := upIPRI.TEIDRange << (8 - upIPRI.TEIDRI)
			teid = uint32(v) << 24

		}
	}

	go RecvProcess(client)

	for i := 0; i < 10; i++ {
		teid++
		sequenceNumber++
		seid++
		time.Sleep(2 * time.Second)

		var pfcpSessionEstablishmentRequest *msg.PFCPSessionEstablishmentRequest
		var err error

		m := urr.NewMeasurementMethod(true, false, false)
		r := urr.NewReportingTriggers(false, false, true, false, false, false, false, false, false, false, false, false, false, false)
		//3600*10s for Time Threshold
		createURR, err := usage_report.NewCreateURR(1, m, r, 0, 36000)

		if ftup {
			fteid, _ := setting.Assign_tunnelID(nil, 0)
			pfcpSessionEstablishmentRequest, err = session.CreateSession(seid, sequenceNumber, nodeIP, seid, 1, 1, 0, fteid, 2, 1, nil, createURR, 1)
		} else {
			fteid, _ := setting.Assign_tunnelID(upIPRI.IPv4Address, teid)
			log.WithFields(log.Fields{"Ftied V4": upIPRI.IPv4Address}).Info("FTEID IPv4 address")

			pfcpSessionEstablishmentRequest, err = session.CreateSession(seid, sequenceNumber, nodeIP, seid, 1, 1, 0, fteid, 2, 1, upIPRI.NetworkInstance, createURR, 1)
		}
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
		seSEIDSN.Inc(seid, sequenceNumber)
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
			var smr *msg.PFCPSessionModificationRequest
			var err error
			seid := cpSEIDDPSEID.Value(srr.SRequest.GetHeader().SessionEndpointIdentifier)
			if seid == 0 {
				seid = srr.SRequest.GetHeader().SessionEndpointIdentifier
			}

			if ftup {
				smr, err = session.ModifySession(seid, sequenceNumber, 2, 2, ie.Core, ueIPAddress, 0, rIPAddress, uint8(ie.FORW), ie.Access, nil)
			} else {
				smr, err = session.ModifySession(seid, sequenceNumber, 2, 2, ie.Core, ueIPAddress, rteid, rIPAddress, uint8(ie.FORW), ie.Access, upIPRI.NetworkInstance)
			}

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
			seid := cpSEIDDPSEID.Value(srr.SRequest.GetHeader().SessionEndpointIdentifier)
			if seid == 0 {
				fmt.Printf("SEID before delete %d\n", srr.SRequest.GetHeader().SessionEndpointIdentifier)
				seid = srr.SRequest.GetHeader().SessionEndpointIdentifier
			}
			pfcpHeader := msg.NewPFCPHeader(1, false, true, msg.SessionDeletionRequestType, 12, seid, sequenceNumber, 0)
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

	return nil

}

func main() {
	app := cli.NewApp()
	app.Name = "N4-Go-Client"
	app.Usage = "N4-Go-Client"
	app.Action = run
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "localIP,l",
			Usage: "localIP address",
			Value: "127.0.0.1",
		},
		&cli.StringFlag{
			Name:  "remoteIP,r",
			Usage: "remoteIP address",
			Value: "127.0.0.1",
		},
		&cli.BoolFlag{
			Name:  "clientInit,c",
			Usage: "clientInit association",
		},
	}

	app.Run(os.Args)

}
