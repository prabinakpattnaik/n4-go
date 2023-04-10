package session

import "C"

import (
	"fmt"
	"net"
	"os"
//	"time"

	"encoding/binary"

	"github.com/prabinakpattnaik/n4-go/ie"
	"github.com/prabinakpattnaik/n4-go/ie/bar"
	"github.com/prabinakpattnaik/n4-go/ie/qer"
	"github.com/prabinakpattnaik/n4-go/msg"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	log_1 "log"
	log "github.com/sirupsen/logrus"
)

var (
	precedance  = 10
	oHR         = []byte{0x00} //GTP-U/UDP/IP4
	UEIPaddress = ""

        WarningLogger *log_1.Logger
        InfoLogger    *log_1.Logger
        ErrorLogger   *log_1.Logger
)

func ProcessPFCPSessionEstablishmentResponse(m *msg.PFCPMessage) ([]byte, error) {

	return nil, nil
}

func CreateSession(sei uint64, sn uint32, nodeIP net.IP, seid uint64, pdrid uint16, farid uint32, sourceinterface ie.InterfaceValue, fteid *ie.FTEID, aa uint8, destionationinterface ie.InterfaceValue, ni string, c *ie.InformationElement, urrid uint32, createQER *qer.CreateQER, qerid uint32, ueipAddress net.IP, teid uint32, remoteIP net.IP, interval uint32, tos uint8) (*msg.PFCPSessionEstablishmentRequest, error) {
	//TODO nodeIP is IPv4 address.
	// Need to change when accomadating FQDN
	// SN incremental (request and response has same value)
	// SEID in increment for each session, set by sending entity. Each session, sending side uses SEID X and receiving SEID Y)
	//
	// Error: Session context not found

	nodeID := []byte{0x00}
	nodeID = append(nodeID, nodeIP.To4()...)
	nodeIDIE := ie.NewInformationElement(
		ie.IENodeID, //IEcode
		0,           //EntrepriseID
		dt.OctetString(nodeID),
	)
	length := ie.IEBasicHeaderSize + nodeIDIE.Len()

	fseid := ie.NewFSEID(true, false, seid, nodeIP, nil)
	bb, err := fseid.Serialize()
	if err != nil {
		return nil, err
	}
	cpfseidIE := ie.NewInformationElement(
		ie.IEFSEID,
		0,
		dt.OctetString(bb),
	)
	length += ie.IEBasicHeaderSize + cpfseidIE.Len()

	si := ie.NewInformationElement(
		ie.IESourceInterface,
		0,
		dt.OctetString(sourceinterface),
	)

	var pdi *ie.PDIWithIE

	if fteid != nil {
		bb, err = fteid.Serialize()
		if err != nil {
			return nil, err
		}
		fteidIE := ie.NewInformationElement(
			ie.IEFTEID,
			0,
			dt.OctetString(bb),
		)

		b_pdr1 := []byte(ni)
	        l := len(b_pdr1)
		b1 := make([]byte, l+1)
	        b1[0] = byte(l)
	        copy(b1[1:], b_pdr1)
	        networkInstance := ie.NewInformationElement(
	                ie.IENetworkInstance,
	                0,
	                dt.OctetString(b1),
	        )

		ueIPAddress := ie.NewUEIPAddress(false, true, false, false, ueipAddress, nil, 0)
		bb, err := ueIPAddress.Serialize()
		if err != nil {
			return nil, err
		}

		ueIPAddressIE := ie.NewInformationElement(
			ie.IEUEIPaddress,
			0,
			dt.OctetString(bb),
		)
		flowdescription := []byte("permit in ip from any to assigned")
                sdffilter := ie.SDFFilter{
			FD:                      true,
                        TTC:                     true,
                        SPI:                     false,
                        FL:                      false,
                        BID:                     false,
                        LengthOfFlowDescription: uint16(len(flowdescription)),
                        FlowDescription:         flowdescription,
                        ToSTrafficClass:         []byte{tos, 0},
		}
		buf, err := sdffilter.Serialize()
		if err != nil {
                        return nil, err
                }
                sdffilterIE := ie.NewInformationElement(
                        ie.IESDFFilter,
                        0,
                        dt.OctetString(buf),
                )

		pdi = ie.NewPDI(&si, &fteidIE, &networkInstance, &ueIPAddressIE, nil, &sdffilterIE, nil, nil, nil, nil, nil, nil, nil)
	} else {
		pdi = ie.NewPDI(&si, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	}
	bb, err = pdi.Serialize()
	if err != nil {
		return nil, err
	}

	pdiIE := ie.NewInformationElement(
		ie.IEPDI,
		0,
		dt.OctetString(bb),
	)

	d := make([]byte, 2)
	binary.BigEndian.PutUint16(d, pdrid)
	pdrIDIE := ie.NewInformationElement(
		ie.IEPDRID,
		0,
		dt.OctetString(d),
	)

	precedenceIE := ie.NewInformationElement(
		ie.IEPrecedence,
		0,
		dt.Unsigned32(precedance),
	)

	outerHeaderRemovalIE := ie.NewInformationElement(
		ie.IEOuterHeaderRemoval,
		0,
		dt.OctetString(oHR),
	)

	farIDIE := ie.NewInformationElement(
		ie.IEFARID,
		0,
		dt.Unsigned32(farid),
	)

	urrIDIE := ie.NewInformationElement(
		ie.IEURRID,
		0,
		dt.Unsigned32(urrid),
	)
	qerIDIE := ie.NewInformationElement(
                ie.IEQERID,
                0,
                dt.Unsigned32(qerid),
        )

	createPDR := ie.NewCreatePDR(&pdrIDIE, &precedenceIE, &pdiIE, &outerHeaderRemovalIE, &farIDIE, &urrIDIE, &qerIDIE, nil)

	bb, err = createPDR.Serialize()
	if err != nil {
		return nil, err
	}
        var createPDRIE ie.InformationElements
	createPDREle := ie.NewInformationElement(
		ie.IECreatePDR,
		0,
		dt.OctetString(bb),
	)
        createPDRIE = append(createPDRIE, createPDREle)
	length = length + ie.IEBasicHeaderSize + createPDREle.Len()



//for PDR 2
//fill pdr 2

	var pdrid1 uint16 = 2
        d1 := make([]byte, 2)
        binary.BigEndian.PutUint16(d1, pdrid1)
        pdrIDIE1 := ie.NewInformationElement(
                ie.IEPDRID,
                0,
                dt.OctetString(d1),
        )

	var farid1 uint32 = 2

        farIDIE1 := ie.NewInformationElement(
                ie.IEFARID,
                0,
                dt.Unsigned32(farid1),
        )

/*	var urrid1 uint32 = 2
        urrIDIE1 := ie.NewInformationElement(
                ie.IEURRID,
                0,
                dt.Unsigned32(urrid1),
        )
*/
	var pdi1 *ie.PDIWithIE
	var sourceinterface1 ie.InterfaceValue = 2
        si1 := ie.NewInformationElement(
                ie.IESourceInterface,
                0,
                dt.OctetString(sourceinterface1),
        )

        if fteid != nil {
/*                bb, err = fteid.Serialize()
                if err != nil {
                        return nil, err
                }
                fteidIE := ie.NewInformationElement(
                        ie.IEFTEID,
                        0,
                        dt.OctetString(bb),
                )*/

                b_pdr1 := []byte(ni)
                l := len(b_pdr1)
                b1 := make([]byte, l+1)
                b1[0] = byte(l)
                copy(b1[1:], b_pdr1)
                networkInstance := ie.NewInformationElement(
                        ie.IENetworkInstance,
                        0,
                        dt.OctetString(b1),
                )

                ueIPAddress := ie.NewUEIPAddress(false, true, true, false, ueipAddress, nil, 0)
                bb, err := ueIPAddress.Serialize()
                if err != nil {
                        return nil, err
                }

                ueIPAddressIE := ie.NewInformationElement(
                        ie.IEUEIPaddress,
                        0,
                        dt.OctetString(bb),
                )
                flowdescription := []byte("permit out ip from any to assigned")
                sdffilter := ie.SDFFilter{
                        FD:                      true,
                        TTC:                     true,
                        SPI:                     false,
                        FL:                      false,
                        BID:                     false,
                        LengthOfFlowDescription: uint16(len(flowdescription)),
                        FlowDescription:         flowdescription,
                        ToSTrafficClass:         []byte{tos, 0},
                }
                buf, err := sdffilter.Serialize()
                if err != nil {
                        return nil, err
                }
                sdffilterIE := ie.NewInformationElement(
                        ie.IESDFFilter,
                        0,
                        dt.OctetString(buf),
                )

                pdi1 = ie.NewPDI(&si1, nil, &networkInstance, &ueIPAddressIE, nil, &sdffilterIE, nil, nil, nil, nil, nil, nil, nil)
        } else {
                pdi1 = ie.NewPDI(&si1, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
        }


        bb, err = pdi1.Serialize()
        if err != nil {
                return nil, err
        }

        pdiIE1 := ie.NewInformationElement(
                ie.IEPDI,
                0,
                dt.OctetString(bb),
        )

	createPDR1 := ie.NewCreatePDR(&pdrIDIE1, &precedenceIE, &pdiIE1, nil, &farIDIE1, &urrIDIE, &qerIDIE, nil)
        bb, err = createPDR1.Serialize()
        if err != nil {
                return nil, err
        }
        createPDREle1 := ie.NewInformationElement(
                 ie.IECreatePDR,
                 0,
                 dt.OctetString(bb),
        )
	createPDRIE = append(createPDRIE, createPDREle1)
	length = length + ie.IEBasicHeaderSize + createPDREle1.Len()




	applyAction := dt.OctetString([]byte{aa})
	applyActionIE := ie.NewInformationElement(
		ie.IEApplyAction,
		0,
		applyAction,
	)

	destionationInterfaceIE := ie.NewInformationElement(
		ie.IEDestinationInterface,
		0,
		dt.OctetString(destionationinterface),
	)
	b := []byte(ni)
	l := len(b)
	b1 := make([]byte, l+1)
	b1[0] = byte(l)
	copy(b1[1:], b)
	networkInstance := ie.NewInformationElement(
		ie.IENetworkInstance,
		0,
		dt.OctetString(b1),
	)

	fp := ie.NewForwardingParameters(&destionationInterfaceIE, &networkInstance, nil, nil, nil, nil, nil, nil, nil)
	bb, err = fp.Serialize()
	if err != nil {
		return nil, err
	}

	forwardingParametersIE := ie.NewInformationElement(
		ie.IEForwardingParameters,
		0,
		dt.OctetString(bb),
	)

	createFAR := ie.NewCreateFAR(&farIDIE, &applyActionIE, &forwardingParametersIE, nil, nil)
	bb, err = createFAR.Serialize()
	if err != nil {
		return nil, err
	}
        var createFARIE ie.InformationElements
	createFAREle := ie.NewInformationElement(
		ie.IECreateFAR,
		0,
		dt.OctetString(bb),
	)
        createFARIE = append(createFARIE, createFAREle)
	length = length + ie.IEBasicHeaderSize + createFAREle.Len()

//FAR2
	var destionationinterface1 =  ie.Access
        destionationInterfaceIE1 := ie.NewInformationElement(
                ie.IEDestinationInterface,
                0,
                dt.OctetString(destionationinterface1),
        )
	ohcd := uint8(1)
	ohc := ie.NewOuterHeaderCreation(ohcd, teid, remoteIP, nil, 0)
	b, err = ohc.Serialize()
	if err != nil {
		return nil, err
	}
	ohcIE := ie.NewInformationElement(
		ie.IEOuterHeaderCreation,
		0,
		dt.OctetString(b),
	)
        b_nwi1 := []byte("access")
        l_1 := len(b_nwi1)
        b1_1 := make([]byte, l_1+1)
        b1_1[0] = byte(l_1)
        copy(b1_1[1:], b_nwi1)
        networkInstance1 := ie.NewInformationElement(
                ie.IENetworkInstance,
                0,
                dt.OctetString(b1_1),
        )

        fp_1 := ie.NewForwardingParameters(&destionationInterfaceIE1, &networkInstance1, nil, &ohcIE, nil, nil, nil, nil, nil)
        bb, err = fp_1.Serialize()
        if err != nil {
                return nil, err
        }

        forwardingParametersIE_1 := ie.NewInformationElement(
                ie.IEForwardingParameters,
                0,
                dt.OctetString(bb),
        )

        createFAR1 := ie.NewCreateFAR(&farIDIE1, &applyActionIE, &forwardingParametersIE_1, nil, nil)
        bb, err = createFAR1.Serialize()
        if err != nil {
                return nil, err
        }

        createFAREle1 := ie.NewInformationElement(
                ie.IECreateFAR,
                0,
                dt.OctetString(bb),
        )
	createFARIE = append(createFARIE, createFAREle1)
	length = length + ie.IEBasicHeaderSize + createFAREle1.Len()


	if c != nil {
		length += ie.IEBasicHeaderSize + c.Len()
	}

	var createQERIE ie.InformationElement
	if createQER != nil {
		bb, err = createQER.Serialize()
		if err != nil {
			return nil, err
		}
		createQERIE = ie.NewInformationElement(
			ie.IECreateQER,
			0,
			dt.OctetString(bb),
		)
		length = length + ie.IEBasicHeaderSize + createQERIE.Len()
	}
        // URR IE Set
	m := ie.NewMeasurementMethod(true, true, false)
	mm, err := m.Serialize()
	if err != nil {
		return nil, err
	}
	MeasurementMethodIE := ie.NewInformationElement(
		ie.IEMeasurementMethod,
		0,
		dt.OctetString(mm),
	)
	ReportingTriggers := ie.NewReportingTriggers(true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false)
	rt, err := ReportingTriggers.Serialize()
	if err != nil {
                return nil, err
        }
        ReportingTriggersIE := ie.NewInformationElement(
		ie.IEReportingTriggers,
		0,
		dt.OctetString(rt),
	)
        MeasurementPeriodIE := ie.NewInformationElement(
		ie.IEMeasurementPeriod,
		0,
		dt.Unsigned32(interval),
	)
	createURR := ie.NewCreateURR(&urrIDIE, &MeasurementMethodIE, &ReportingTriggersIE, &MeasurementPeriodIE, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	bb, err = createURR.Serialize()
	if err != nil {
                return nil, err
        }
	var createURRie ie.InformationElement
        createURRie = ie.NewInformationElement(
		ie.IECreateURR,
		0,
		dt.OctetString(bb),
	)
	length = length + ie.IEBasicHeaderSize + createURRie.Len()
	//user plane inactivity Timer

	upInactivityTimerIE := ie.NewInformationElement(
		ie.IEUserPlaneInactivityTimer,
		0,
		dt.Unsigned32(300),
	)
	length = length + ie.IEBasicHeaderSize + upInactivityTimerIE.Len()

	pfcpHeader := msg.NewPFCPHeader(1, false, true, msg.SessionEstablishmentRequestType, length+12, sei, sn, 0)
	if createQERIE.Type == ie.IEReserved {
		pfcpSessionEstablishmentRequest := msg.NewPFCPSessionEstablishmentRequest(pfcpHeader, &nodeIDIE, &cpfseidIE, &createPDRIE, &createFARIE, c, nil, nil, nil, nil, &upInactivityTimerIE, nil, nil)
		return &pfcpSessionEstablishmentRequest, nil
	} else {
		pfcpSessionEstablishmentRequest := msg.NewPFCPSessionEstablishmentRequest(pfcpHeader, &nodeIDIE, &cpfseidIE, &createPDRIE, &createFARIE, &createURRie, &createQERIE, nil, nil, nil, &upInactivityTimerIE, nil, nil)
		return &pfcpSessionEstablishmentRequest, nil
	}
	return nil, nil

}

func ModifySession(sei uint64, sn uint32, pdrid uint16, farid uint32, sourceinterface ie.InterfaceValue, ueipAddress net.IP, teid uint32, remoteIP net.IP, aa ie.ApplyActionValue, dInterface ie.InterfaceValue, ni []byte, createBAR *bar.CreateBAR) (*msg.PFCPSessionModificationRequest, error) {

        file, err := os.OpenFile("/home/vagrant/magma/lte/gateway/log_2s.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
        if err != nil {
            log_1.Fatal(err)
        }

        InfoLogger = log_1.New(file, "INFO: ", log_1.Ldate|log_1.Ltime|log_1.Lshortfile)
        WarningLogger = log_1.New(file, "WARNING: ", log_1.Ldate|log_1.Ltime|log_1.Lshortfile)
        ErrorLogger = log_1.New(file, "ERROR: ", log_1.Ldate|log_1.Ltime|log_1.Lshortfile)

	d := make([]byte, 2)
	binary.BigEndian.PutUint16(d, pdrid)
	pdrIDIE := ie.NewInformationElement(
		ie.IEPDRID,
		0,
		dt.OctetString(d),
	)

	precedence := ie.NewInformationElement(
		ie.IEPrecedence,
		0,
		dt.Unsigned32(precedance),
	)

	farIDIE := ie.NewInformationElement(
		ie.IEFARID,
		0,
		dt.Unsigned32(farid),
	)
	si := ie.NewInformationElement(
		ie.IESourceInterface,
		0,
		dt.OctetString(byte(sourceinterface)),
	)
	//TODO ueipAddress is IPv4 address
	ueIPAddress := ie.NewUEIPAddress(false, true, true, false, ueipAddress, nil, 0)
	bb, err := ueIPAddress.Serialize()
	if err != nil {
		return nil, err
	}

	ueIPAddressIE := ie.NewInformationElement(
		ie.IEUEIPaddress,
		0,
		dt.OctetString(bb),
	)

	b := []byte("internet")
	l := len(b)
	b1 := make([]byte, l+1)
	b1[0] = byte(l)
	copy(b1[1:], b)

	var networkInstance ie.InformationElement
	networkInstance = ie.NewInformationElement(
		ie.IENetworkInstance,
		0,
		dt.OctetString(b1),
	)

	pdi := ie.NewPDI(&si, nil, &networkInstance, &ueIPAddressIE, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	pdiB, err := pdi.Serialize()
	if err != nil {
		return nil, err
	}

	pdiIE := ie.NewInformationElement(
		ie.IEPDI,
		0,
		dt.OctetString(pdiB),
	)

	updatePDR := ie.NewUpdatePDR(&pdrIDIE, nil, &precedence, &pdiIE, &farIDIE, nil, nil, nil, nil)
	//createPDR := ie.UpdatePDR(&pdrIDIE, &precedence, &pdiIE, nil, &farIDIE, nil, nil, nil)
	b, err = updatePDR.Serialize()
	log.WithFields(log.Fields{"data": b}).Info("UpdatePDR")
	if err != nil {
		return nil, err

	}
	var updatePDREle ie.InformationElement
	updatePDREle = ie.NewInformationElement(
		ie.IEUpdatePDR,
		0,              //EntrepriseID
		dt.OctetString(b),
	)

	var length uint16 = 0
/*       for i := 0; i <1; i++ {
	        createPDRIE = append(createPDRIE, createPDREle)
		length = ie.IEBasicHeaderSize + createPDREle.Len()
        }*/
	length = ie.IEBasicHeaderSize + updatePDREle.Len()
	var aaValue uint8
	switch aa {

	case ie.DROP:
		aaValue = 1
	case ie.FORW:
		aaValue = 2
	case ie.BUFF:
		aaValue = 4
	case ie.NOCP:
		aaValue = 8
	case ie.DUPL:
		aaValue = 16
	default:
		return nil, fmt.Errorf("Not valid Apply Action")
	}

	aaIE := ie.NewInformationElement(
		ie.IEApplyAction,
		0,
		dt.OctetString(byte(aaValue)),
	)

	desIE := ie.NewInformationElement(
		ie.IEDestinationInterface,
		0,
		dt.OctetString(byte(dInterface)),
	)

	networkInstance = ie.NewInformationElement(
		ie.IENetworkInstance,
		0,
		dt.OctetString(ni),
	)
	var updateFAR *ie.UpdateFARWithIE
	if aa == ie.FORW {
		InfoLogger.Println("Under outer header creation")
		ohcd := uint8(1)
		ohc := ie.NewOuterHeaderCreation(ohcd, 100, remoteIP, nil, 0)
		b, err = ohc.Serialize()
		if err != nil {
			return nil, err
		}
		ohcIE := ie.NewInformationElement(
			ie.IEOuterHeaderCreation,
			0,
			dt.OctetString(b),
		)
		b_nwi1 := []byte("access")
	        l_1 := len(b_nwi1)
	        b1_1 := make([]byte, l_1+1)
	        b1_1[0] = byte(l_1)
	        copy(b1_1[1:], b_nwi1)
	        networkInstance1 := ie.NewInformationElement(
	                ie.IENetworkInstance,
	                0,
	                dt.OctetString(b1_1),
	        )
		fp := ie.NewUpdateForwardingParameters(&desIE, &networkInstance1, nil, &ohcIE, nil, nil, nil, nil, nil)
		bb, err = fp.Serialize()
		if err != nil {
			return nil, err
		}

		fpIE := ie.NewInformationElement(
			ie.IEUpdateForwardingParamets,
			0,
			dt.OctetString(bb),
		)

		updateFAR = ie.NewUpdateFAR(&farIDIE, &aaIE, &fpIE, nil, nil)
		//createFAR = ie.UpdateFAR(&farIDIE, &aaIE, &fpIE, nil, nil)

	} else if aa == ie.BUFF {
		updateFAR = ie.NewUpdateFAR(&farIDIE, &aaIE, nil, nil, createBAR.BARID)

	}else if aa == ie.DROP {
		updateFAR = ie.NewUpdateFAR(&farIDIE, &aaIE, nil, nil, nil)
	}

	bCreateFAR, err := updateFAR.Serialize()
	if err != nil {
		return nil, err

	}
        var updateFAREle ie.InformationElement
	updateFAREle = ie.NewInformationElement(
		ie.IEUpdateFAR,
		0,
		dt.OctetString(bCreateFAR),
	)
	length = length + ie.IEBasicHeaderSize + updateFAREle.Len()
	var smr msg.PFCPSessionModificationRequest
	if aa == ie.FORW || aa == ie.DROP {
		InfoLogger.Println("For msg.NewPFCPSessionModificationRequest")
		pfcpHeader := msg.NewPFCPHeader(1, false, true, msg.SessionModificationRequestType, length+12, sei, sn, 0)
		//smr = msg.NewPFCPSessionModificationRequest(pfcpHeader, nil, nil, nil, nil, nil, nil, nil, &createPDRIE, &createFARIE, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	        log.WithFields(log.Fields{"data": updatePDREle}).Info("UpdatePDRIE")
		log.WithFields(log.Fields{"data": updateFAREle}).Info("updateFAREle")
		smr = msg.NewPFCPSessionModificationRequest(pfcpHeader, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &updatePDREle, &updateFAREle, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	} else if aa == ie.BUFF {
		bcreateBAR, err := createBAR.Serialize()
		if err != nil {
			return nil, err

		}
		fmt.Printf("[%x]\n", bcreateBAR)
		createBARIE := ie.NewInformationElement(
			ie.IECreateBAR,
			0,
			dt.OctetString(bcreateBAR),
		)
		length = length + ie.IEBasicHeaderSize + createBARIE.Len()
		pfcpHeader := msg.NewPFCPHeader(1, false, true, msg.SessionModificationRequestType, length+12, sei, sn, 0)
		//smr = msg.NewPFCPSessionModificationRequest(pfcpHeader, nil, nil, nil, nil, nil, nil, nil, &createPDRIE, &createFARIE, nil, nil, &createBARIE, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		smr = msg.NewPFCPSessionModificationRequest(pfcpHeader, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &createBARIE, nil, &updatePDREle, &updateFAREle, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	}
	return &smr, nil

}
