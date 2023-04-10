package main
import "C" 

import (
	"fmt"
	"net"
	"os"
	"time"
	//"github.com/prabinakpattnaik/n4-go/ie/qer"
	"github.com/prabinakpattnaik/n4-go/ie/qer"
	"github.com/prabinakpattnaik/n4-go/msg"

	"github.com/prabinakpattnaik/n4-go/util/se"

        setting "github.com/prabinakpattnaik/n4-go/client/internal/helper"
	"github.com/prabinakpattnaik/n4-go/client/internal/server_wrap"
	"github.com/prabinakpattnaik/n4-go/client/internal/session"
	/*setting "client/internal/helper"
        "client/internal/server_wrap"
        "client/internal/session"*/

	"github.com/prabinakpattnaik/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	log_1 "log"
	"context"

        "github.com/golang/glog"
        lte_protos "magma/lte/cloud/go/protos"
        "magma/orc8r/lib/go/merrors"
	"google.golang.org/grpc"

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
	NetworkInstance                   = "internet"
	seSEIDSN                          = session.SESEIDSNEntity{M: make(map[uint64]uint32)}
	cpSEIDDPSEID                      = &se.CPSEIDDPSEIDEntity{M: make(map[uint64]uint64)}
	seSEIDCAUSE                       = session.SEIDCauseEntity{M: make(map[uint64]uint8)}
	seSEIDsessionInfo                 = session.SEIDSessionInfo{M: make(map[uint64]session.SessionInfo)}
	counter                           = 1

	WarningLogger *log_1.Logger
	InfoLogger    *log_1.Logger
	ErrorLogger   *log_1.Logger
	localNode_Id                string
	remoteNode_Id               string
	n3_interface_Id             string
        usagereport_interv          uint32 = 10
)

// Client implements a PFCP client
type Client struct {
	LocalAddr  net.Addr
	RemoteAddr *net.UDPAddr
	Conn       *net.UDPConn
}
//var g_client *Client
var client *Client
var g_ftup bool
var g_upIPRI ie.UPIPResourceInformation
var Epoch uint64
// Acceptance in a response
const (
        CauseRequestAccepted uint8 = 1
)
const (
	Sessiond_grpc_address = "localhost:50065"
)
// Rejection in a response
const (
        CauseRequestRejected uint8 = iota + 64
        CauseSessionContextNotFound
        CauseMandatoryIeMissing
        CauseConditionalIeMissing
        CauseInvalidLength
        CauseMandatoryIeIncorrect
        CauseInvalidForwardingPolicy
        CauseInvalidFTeidAllocationOption
        CauseNoEstablishedPfcpAssociation
        CauseRuleCreationModificationFailure
        CausePfcpEntityInCongestion
        CauseNoResourcesAvailable
        CauseServiceNotSupported
        CauseSystemFailure
)

type Cause struct {
        CauseValue uint8
}
/*
func getEpochTime() int64 {
    return time.Now().Unix()
}*/
func getClient() (lte_protos.LocalSessionManagerClient, error) {
	conn, err := grpc.Dial(Sessiond_grpc_address, grpc.WithInsecure(), grpc.WithBlock())
        if err != nil {
                initErr := merrors.NewInitError(err, "sessiond")
                glog.Error(initErr)
                return nil, initErr
        }
        return lte_protos.NewLocalSessionManagerClient(conn), nil
}

func ReportsRuleStats(ctx context.Context, record *lte_protos.RuleRecord) error {
        client, err := getClient()
        if err != nil {
                return err
        }
        _, err = client.ReportRuleStats(
                ctx,
                &lte_protos.RuleRecordTable{
			Records: []*lte_protos.RuleRecord{record},
                        Epoch: Epoch,
                        UpdateRuleVersions: false,
                },
        )
        if err != nil {
                return err
        }

        return nil
}

// NewClient returns a Client with default settings
func NewClient(lAddress string, rAddress string) *Client {
	remoteIPv4address := net.ParseIP(rAddress)
	raddr := fmt.Sprintf("%s:%d", remoteIPv4address, udpport)
	dst, err := net.ResolveUDPAddr("udp", raddr)

	if err != nil {
		log.WithError(err).Error("resolveUDP dst Addr err ")
	}
	localIPv4address := net.ParseIP(lAddress)
	laddr := fmt.Sprintf("%s:%d", localIPv4address, udpport)
	src, err := net.ResolveUDPAddr("udp", laddr)

	if err != nil {
		log.WithError(err).Error("resolveUDP src Addr err ")
	}

	conn, err := net.DialUDP("udp4", src, dst)
	InfoLogger.Println("value of conn is:", conn)
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
func usagereportpush(uplink_bytes uint64, downlink_bytes uint64, feid uint64) {
	seinfo := seSEIDsessionInfo.Value(feid)
	rulerecord := &lte_protos.RuleRecord{
                Sid:        seinfo.Sid,
                RuleId:     seinfo.Rule_id,
                BytesTx:     uplink_bytes,
                BytesRx:     downlink_bytes,
                UeIpv4:      seinfo.Ue_ipv4,
                DroppedTx:   0,
                DroppedRx:   0,
                RuleVersion: 1,
                Teid:         seinfo.Teid,
         }
	 ctx := context.Background()
	 err := ReportsRuleStats(ctx, rulerecord)
	 if err != nil {
		 log.WithError(err).Info("Error in ReportRuleStats")
	 }
}

func gbr_5QI_dscp(QI_value uint8) uint8 {
	switch QI_value {
        case 1:
		return 46
        case 2:
                return 36
        case 3:
                return 10
        case 4:
                return 28
        case 65:
                return 46
        case 66:
                return 46
        case 67:
                return 34
        case 71:
        case 72:
        case 73:
        case 74:
        case 76:
                return 28
        default:
                return 0
        }
	return 0
}

func non_gbr_5QI_dscp(QI_value uint8) uint8 {
	switch QI_value {
        case 5:
		return 46
	case 6:
		return 10
	case 7:
		return 18
	case 8:
		return 18
	case 9:
		return 0
	case 69:
		return 34
	case 70:
		return 18
	case 79:
		return 18
	case 80:
		return 68
	default:
		return 0
	}
	return 0
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

				        if pfcpSessionEstablishmentResponse.Cause != nil {
						c, _ := pfcpSessionEstablishmentResponse.Cause.Serialize()
					        CauseValue := c[4]
					        seSEIDCAUSE.Inc(fseid.SEID, CauseValue)
					        if CauseValue == CauseRequestAccepted {
							log.WithFields(log.Fields{"cause": CauseValue}).Info("Received PFCP Session Establishment Accepted Response")
                                                } else {
							log.WithFields(log.Fields{"cause": CauseValue}).Info("Received PFCP Session Establishment Not Accepted Response")
					        }
					}
                                       // TODO: set appropriate 5GSM cause according to PFCP cause value
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
				/*
					r := ie.NewInformationElement(
						ie.IERecoveryTimestamp,
						0,
						dt.Time(time.Now()),
					)
				*/
				h := pfcpHeartbeat.GetHeader()
				h.MessageType = msg.HeartbeatResponseType
				//h.MessageLength = msg.PFCPBasicMessageSize + ie.IEBasicHeaderSize + r.Len()
				heartbeat := msg.NewHeartbeat(h, pfcpHeartbeat.RecoveryTimeStamp)
				b, _ := heartbeat.Serialize()
				c.Write(b)
				continue
			}

			//pfcpSessionReportRequest
			if pfcpMessage.Header.MessageType == msg.SessionReportRequestType {
				log.WithFields(log.Fields{"data": rb}).Info("received PFCP Session Report Request")
				b, err := msg.ProcessPFCPSessionReportRequest(pfcpMessage, cpSEIDDPSEID)
				if err != nil {
					log.WithError(err).Error("Process error in PFCP Session Report Request")
				}
				c.Write(b)
				pfcpMessage1, err := msg.MessageFromBytes(rb)
				bc, seid, err := msg.ProcessPFCPUsageReport_seid(pfcpMessage1, cpSEIDDPSEID)
				if err != nil {
                                        log.WithError(err).Error("Process error in PFCP usage Report")
                                }

				bd := bc.VolumeMeasurement
				Total_volume  := bd.TotalVolume
				downlink_volume := bd.DownlinkVolume
				uplink_volume := bd.UplinkVolume
				log.WithFields(log.Fields{"seid": seid}).Info("seid")
				log.WithFields(log.Fields{"Total_volume": Total_volume}).Info("Total_volume")
				log.WithFields(log.Fields{"uplink_volume": uplink_volume}).Info("uplink_volume")
				log.WithFields(log.Fields{"downlink_volume": downlink_volume}).Info("downlink_volume")
				usagereportpush(uplink_volume,  downlink_volume, seid)
				continue
			}

			log.WithFields(log.Fields{"data [%x]": rb}).Info("something went")
		}
	}

}

//export StartPfcpProcess
func StartPfcpProcess(local_ip *C.char, remote_ip *C.char, n3_interface *C.char, usagereport_interval uint32, epoch uint64) int64 {

	file, err := os.OpenFile("/tmp/log_1s.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
	    log_1.Fatal(err)
	}

	InfoLogger = log_1.New(file, "INFO: ", log_1.Ldate|log_1.Ltime|log_1.Lshortfile)
	WarningLogger = log_1.New(file, "WARNING: ", log_1.Ldate|log_1.Ltime|log_1.Lshortfile)
	ErrorLogger = log_1.New(file, "ERROR: ", log_1.Ldate|log_1.Ltime|log_1.Lshortfile)
        localNode_Id = C.GoString(local_ip)
	remoteNode_Id = C.GoString(remote_ip)
	n3_interface_Id = C.GoString(n3_interface)
	usagereport_interv = usagereport_interval
	Epoch = epoch
        app := cli.NewApp()
        app.Name = "N4-Go-Client"
        app.Usage = "N4-Go-Client"
        app.Action = run
        app.Flags = []cli.Flag{
                &cli.StringFlag{
                        Name:  "localIP,l",
                        Usage: "localIP address",
                        Value: C.GoString(local_ip),
                },
                &cli.StringFlag{
                        Name:  "remoteIP,r",
                        Usage: "remoteIP address",
                        Value: C.GoString(remote_ip),
                },
                &cli.StringFlag{
                        Name:  "ueIP,ueip",
                        Usage: "ueIP address",
                        Value: "10.10.1.2",
                },
                &cli.BoolFlag{
                        Name:  "clientInit,c",
                        Usage: "clientInit association",
                },
        }

        app.Run(os.Args)

	return 0
}

func run(c *cli.Context) error {
	lIPv4address := c.String("localIP")
	nodeIP := net.ParseIP(lIPv4address)
	controlPlaneNodeID := []byte{0x0}
	controlPlaneNodeID = append(controlPlaneNodeID, nodeIP.To4()...)

	rIPv4address := c.String("remoteIP")
	client = NewClient(lIPv4address, rIPv4address)
	var upIPRI ie.UPIPResourceInformation
	var upIPRIs []ie.UPIPResourceInformation
	upIPRIC := make(chan []ie.UPIPResourceInformation)
	ftupC := make(chan bool)
	var ftup bool

	if c.Bool("clientInit") {
		go server_wrap.Run(lIPv4address, udpport, controlPlaneNodeID, upIPRIC, ftupC, cpSEIDDPSEID)

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

		}
	}
	g_ftup = ftup

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

	return nil

}
/*func upf_n4_message_create(local_ip *C.char, ue_IP *C.char, n3_ip *C.char,
        in_tei uint32, sourceinterface uint8, aa uint8,
        destionationinterface uint8, apn *C.char, out_tei uint32,
        gnb_ip *C.char, session_state uint8, session_id uint64,
        qfi uint8, apn_mbr_ul uint64, apn_mbr_dl uint64) uint8 {
*/
//export upf_n4_message_create
func upf_n4_message_create(ue_IP *C.char,in_tei uint32, sourceinterface uint8, aa uint8,
        destionationinterface uint8, apn *C.char, out_tei uint32,
        gnb_ip *C.char, session_state uint8, qfi uint8, apn_mbr_ul uint64,
	apn_mbr_dl uint64, imsi *C.char, rule_id *C.char) uint8 {
	nodeIP := net.ParseIP(localNode_Id)
	ue_ip := net.ParseIP(C.GoString(ue_IP))
	fteid_ip := net.ParseIP(n3_interface_Id)
	rIPAddress:= net.ParseIP(C.GoString(gnb_ip))
	sn_del := uint32(101)
	session_Id := uint64(in_tei)

        for i := 0; i < counter; i++ {
                sequenceNumber++
                var pfcpSessionEstablishmentRequest *msg.PFCPSessionEstablishmentRequest
                var err error

                fteid, _ := setting.Assign_tunnelID(fteid_ip, in_tei)
                log.WithFields(log.Fields{"Ftied V4": fteid_ip}).Info("FTEID IPv4 address")

                //QER
                gateStatus := qer.NewGateStatus(qer.OPEN, qer.OPEN)
                mbr := qer.NewBR(apn_mbr_ul, apn_mbr_dl)
                gbr := qer.NewBR(80, 80)
                rqi := true
                createQER, err := qer.NewCreateQER(1, 0, gateStatus, mbr, gbr, qfi, rqi)
                if err != nil {
                        log.WithError(err).Error("error in creating CreateQER")
                        continue
                }
		var tos uint8
		tos = non_gbr_5QI_dscp(qfi)
		if session_state == 0 {
			seinfo := seSEIDsessionInfo.Value(session_Id)
			if seinfo.Sid != "" {
				previous_state := seinfo.Session_state
				if previous_state == 2 {
					var smr *msg.PFCPSessionModificationRequest
                                        smr, err = session.ModifySession(session_Id, sequenceNumber+1, 2, 2, ie.SGiN6LAN, ue_ip, 0, rIPAddress, ie.FORW, ie.Access, nil, nil)
                                        if err != nil {
						log.WithError(err).Error("error in pfcpSessionModificationRequest")
                                                continue
                                        }
                                        b, err2 := smr.Serialize()
                                        if err2 != nil {
						log.WithError(err).Error("error in pfcpSessionModificationRequest serialization")
                                                continue
                                        }
					seSEIDCAUSE.Inc(session_Id, 0)
					client.Write(b)
					ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
					defer cancel()
					for seSEIDCAUSE.Value(session_Id) == 0 || ctx.Err() != nil {
						InfoLogger.Println("Waiting pfcp session modification response ")
					}
					sessionRequestResponse := session.SessionRequestResponse{
						SRequest: smr,
                                        }
                                        SessionInfo := session.SessionInfo{
						Sid:     C.GoString(imsi),
                                                Rule_id: C.GoString(rule_id),
                                                Ue_ipv4: C.GoString(ue_IP),
                                                Teid: in_tei,
                                                Session_state: session_state,
                                        }
                                        seSEIDsessionInfo.Inc(session_Id, SessionInfo)
					sessionEntity.Inc(sequenceNumber, sessionRequestResponse)
				}
			} else {
				pfcpSessionEstablishmentRequest, err = session.CreateSession(session_Id, sequenceNumber, nodeIP, session_Id, 1, 1, 0, fteid, aa, ie.SGiN6LAN, C.GoString(apn), nil, 1, createQER, 1, ue_ip, out_tei, rIPAddress, usagereport_interv, tos)
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
                                seSEIDCAUSE.Inc(session_Id, 0)
                                client.Write(b)
			        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			        defer cancel()
			        for seSEIDCAUSE.Value(session_Id) == 0 || ctx.Err() != nil {
					InfoLogger.Println("Waiting for pfcp session establishment response")
                                }
			        SessionInfo := session.SessionInfo{
					Sid:     C.GoString(imsi),
				        Rule_id: C.GoString(rule_id),
				        Ue_ipv4: C.GoString(ue_IP),
				        Teid: in_tei,
				        Session_state: session_state,
			        }
				seSEIDsessionInfo.Inc(session_Id, SessionInfo)
			}
		}else if session_state == 1  {
			srr := sessionEntity.Value(sn_del)
			if srr.SRequest != nil {
				seid_del := cpSEIDDPSEID.Value(srr.SRequest.GetHeader().SessionEndpointIdentifier)
                                if seid_del == 0 {
					fmt.Printf("SEID before delete %d\n", srr.SRequest.GetHeader().SessionEndpointIdentifier)
                                        seid_del = srr.SRequest.GetHeader().SessionEndpointIdentifier
				}
                                pfcpHeader := msg.NewPFCPHeader(1, false, true, msg.SessionDeletionRequestType, 12, session_Id, sequenceNumber, 0)
                                b := pfcpHeader.Serialize()
				seSEIDCAUSE.Inc(session_Id, 0)
                                client.Write(b)
                                log.WithFields(log.Fields{"data": b}).Info("received pfcpSessionDeletionRequest")
                                sn_del++
                                sdr := msg.NewPFCPSessionDeletionRequest(pfcpHeader)
                                sessionRequestResponse := session.SessionRequestResponse{
                                       SRequest: &sdr,
                                }
				for seSEIDCAUSE.Value(session_Id) == 0 {
					InfoLogger.Println("Waiting for pfcp session delete response")
				}
                                sessionEntity.Inc(sequenceNumber, sessionRequestResponse)
                                seidsnEntity.Inc(srr.SRequest.GetHeader().SessionEndpointIdentifier, 3, sequenceNumber)
				seSEIDCAUSE.Delete(session_Id)
				seSEIDsessionInfo.Delete(session_Id)
			}
		}else if session_state == 2{
			 var smr *msg.PFCPSessionModificationRequest
                         smr, err = session.ModifySession(session_Id, sequenceNumber+1, 2, 2, ie.SGiN6LAN, ue_ip, 0, rIPAddress, ie.DROP, ie.Access, nil, nil)
                         if err != nil {
                                log.WithError(err).Error("error in pfcpSessionModificationRequest")
                                continue
                         }
                         b, err2 := smr.Serialize()
                         if err2 != nil {
                                log.WithError(err).Error("error in pfcpSessionModificationRequest serialization")
                                continue
                         }
			 seSEIDCAUSE.Inc(session_Id, 0)
                         client.Write(b)
			 sessionRequestResponse := session.SessionRequestResponse{
                                SRequest: smr,
                        }
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			for seSEIDCAUSE.Value(session_Id) == 0 || ctx.Err() != nil {
				InfoLogger.Println("Waiting for pfcp session modification response")
			}
			SessionInfo := session.SessionInfo{
                                Sid:     C.GoString(imsi),
                                Rule_id: C.GoString(rule_id),
                                Ue_ipv4: C.GoString(ue_IP),
                                Teid: in_tei,
                                Session_state: session_state,
                        }
                        seSEIDsessionInfo.Inc(session_Id, SessionInfo)
			sessionEntity.Inc(sequenceNumber, sessionRequestResponse)

		}else {
		}
        }
	return seSEIDCAUSE.Value(session_Id)
}

func main() {
	//aniket
/*        var y int = 2
	var z string = "Aniket"
        demo(y, z)
	app := cli.NewApp()
	app.Name = "N4-Go-Client"
	app.Usage = "N4-Go-Client"
	app.Action = run
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "localIP,l",
			Usage: "localIP address",
			Value: "172.32.1.1",
		},
		&cli.StringFlag{
			Name:  "remoteIP,r",
			Usage: "remoteIP address",
			Value: "172.32.1.2",
		},
		&cli.BoolFlag{
			Name:  "clientInit,c",
			Usage: "clientInit association",
		},
	}

	app.Run(os.Args)*/

}
