package session

import (
	"fmt"
	"net"

	"encoding/binary"

	"bitbucket.org/sothy5/n4-go/ie"
	"bitbucket.org/sothy5/n4-go/ie/bar"
	"bitbucket.org/sothy5/n4-go/ie/qer"
	"bitbucket.org/sothy5/n4-go/msg"
	dt "github.com/fiorix/go-diameter/diam/datatype"
)

var (
	precedance = 10
	oHR        = []byte{0x06} //GTP-U/UDP/IP
)

func ProcessPFCPSessionEstablishmentResponse(m *msg.PFCPMessage) ([]byte, error) {

	return nil, nil
}

func CreateSession(sei uint64, sn uint32, nodeIP net.IP, seid uint64, pdrid uint16, farid uint32, sourceinterface uint8, fteid *ie.FTEID, aa, destionationinterface uint8, ni []byte, c *ie.InformationElement, urrid uint32, createQER *qer.CreateQER, qerid uint32) (*msg.PFCPSessionEstablishmentRequest, error) {
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
		dt.OctetString([]byte{sourceinterface}),
	)

	var pdi *ie.PDI

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
		var networkInstance ie.InformationElement
		if len(ni) > 0 {
			networkInstance = ie.NewInformationElement(
				ie.IENetworkInstance,
				0,
				dt.OctetString(ni),
			)
		}

		pdi = ie.NewPDI(&si, &fteidIE, &networkInstance, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
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

	createPDRIE := ie.NewInformationElement(
		ie.IECreatePDR,
		0,
		dt.OctetString(bb),
	)
	length = length + ie.IEBasicHeaderSize + createPDRIE.Len()

	applyAction := dt.OctetString([]byte{aa})
	applyActionIE := ie.NewInformationElement(
		ie.IEApplyAction,
		0,
		applyAction,
	)

	destionationInterfaceIE := ie.NewInformationElement(
		ie.IEDestinationInterface,
		0,
		dt.OctetString([]byte{destionationinterface}),
	)

	fp := ie.NewForwardingParameters(&destionationInterfaceIE, nil, nil, nil, nil, nil, nil, nil, nil)
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
	createFARIE := ie.NewInformationElement(
		ie.IECreateFAR,
		0,
		dt.OctetString(bb),
	)
	length = length + ie.IEBasicHeaderSize + createFARIE.Len()
	length += ie.IEBasicHeaderSize + c.Len()

	bb, err = createQER.Serialize()
	if err != nil {
		return nil, err
	}
	createQERIE := ie.NewInformationElement(
		ie.IECreateQER,
		0,
		dt.OctetString(bb),
	)
	length = length + ie.IEBasicHeaderSize + createQERIE.Len()

	pfcpHeader := msg.NewPFCPHeader(1, false, true, msg.SessionEstablishmentRequestType, length+12, sei, sn, 0)
	pfcpSessionEstablishmentRequest := msg.NewPFCPSessionEstablishmentRequest(pfcpHeader, &nodeIDIE, &cpfseidIE, &createPDRIE, &createFARIE, c, &createQERIE, nil, nil, nil, nil, nil, nil)

	return &pfcpSessionEstablishmentRequest, nil

}

func ModifySession(sei uint64, sn uint32, pdrid uint16, farid uint32, sourceinterface ie.InterfaceValue, ueipAddress net.IP, teid uint32, remoteIP net.IP, aa ie.ApplyActionValue, dInterface ie.InterfaceValue, ni []byte, createBAR *bar.CreateBAR) (*msg.PFCPSessionModificationRequest, error) {
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

	var networkInstance ie.InformationElement
	if len(ni) > 0 {
		networkInstance = ie.NewInformationElement(
			ie.IENetworkInstance,
			0,
			dt.OctetString(ni),
		)
	}

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

	createPDR := ie.NewCreatePDR(&pdrIDIE, &precedence, &pdiIE, nil, &farIDIE, nil, nil, nil)
	b, err := createPDR.Serialize()
	if err != nil {
		return nil, err

	}

	createPDRIE := ie.NewInformationElement(
		ie.IECreatePDR, //IEcode
		0,              //EntrepriseID
		dt.OctetString(b),
	)
	length := ie.IEBasicHeaderSize + createPDRIE.Len()
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
	var createFAR *ie.CreateFAR
	if aa == ie.FORW {
		ohcd := uint8(1)
		ohc := ie.NewOuterHeaderCreation(ohcd, teid, ueipAddress, nil, 0)
		b, err = ohc.Serialize()
		if err != nil {
			return nil, err
		}
		ohcIE := ie.NewInformationElement(
			ie.IEOuterHeaderCreation,
			0,
			dt.OctetString(b),
		)
		fp := ie.NewForwardingParameters(&desIE, nil, nil, &ohcIE, nil, nil, nil, nil, nil)
		bb, err = fp.Serialize()
		if err != nil {
			return nil, err
		}

		fpIE := ie.NewInformationElement(
			ie.IEForwardingParameters,
			0,
			dt.OctetString(bb),
		)

		createFAR = ie.NewCreateFAR(&farIDIE, &aaIE, &fpIE, nil, nil)

	} else if aa == ie.BUFF {
		createFAR = ie.NewCreateFAR(&farIDIE, &aaIE, nil, nil, createBAR.BARID)

	}

	bCreateFAR, err := createFAR.Serialize()
	if err != nil {
		return nil, err

	}

	createFARIE := ie.NewInformationElement(
		ie.IECreateFAR,
		0,
		dt.OctetString(bCreateFAR),
	)
	length = length + ie.IEBasicHeaderSize + createFARIE.Len()
	var smr msg.PFCPSessionModificationRequest
	if aa == ie.FORW {
		pfcpHeader := msg.NewPFCPHeader(1, false, true, msg.SessionModificationRequestType, length+12, sei, sn, 0)
		smr = msg.NewPFCPSessionModificationRequest(pfcpHeader, nil, nil, nil, nil, nil, nil, nil, &createPDRIE, &createFARIE, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
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
		smr = msg.NewPFCPSessionModificationRequest(pfcpHeader, nil, nil, nil, nil, nil, nil, nil, &createPDRIE, &createFARIE, nil, nil, &createBARIE, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	}
	return &smr, nil

}
