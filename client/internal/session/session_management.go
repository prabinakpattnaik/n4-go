package session

import (
	"net"

	"encoding/binary"

	"bitbucket.org/sothy5/n4-go/ie"
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

func CreateSession(sei uint64, sn uint32, nodeIP net.IP, seid uint64, pdrid uint16, farid uint32, sourceinterface uint8, fteid *ie.FTEID, aa, destionationinterface uint8) (*msg.PFCPSessionEstablishmentRequest, error) {
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

	bb, err = fteid.Serialize()
	if err != nil {
		return nil, err
	}
	fteidIE := ie.NewInformationElement(
		ie.IEFTEID,
		0,
		dt.OctetString(bb),
	)

	pdi := ie.NewPDI(&si, &fteidIE, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
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

	createPDR := ie.NewCreatePDR(&pdrIDIE, &precedenceIE, &pdiIE, &outerHeaderRemovalIE, &farIDIE, nil, nil, nil)
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

	pfcpHeader := msg.NewPFCPHeader(1, false, true, msg.SessionEstablishmentRequestType, length+12, sei, sn, 0)
	pfcpSessionEstablishmentRequest := msg.NewPFCPSessionEstablishmentRequest(pfcpHeader, &nodeIDIE, &cpfseidIE, &createPDRIE, &createFARIE, nil, nil, nil, nil, nil, nil, nil, nil)

	return &pfcpSessionEstablishmentRequest, nil

}

func ModifySession(sei uint64, sn uint32, pdrid uint16, farid uint32, sourceinterface ie.InterfaceValue, ueipAddress net.IP, teid uint32, remoteIP net.IP, aa uint8, dInterface ie.InterfaceValue) (*msg.PFCPSessionModificationRequest, error) {
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

	pdi := ie.NewPDI(&si, nil, nil, &ueIPAddressIE, nil, nil, nil, nil, nil, nil, nil, nil, nil)
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

	aaIE := ie.NewInformationElement(
		ie.IEApplyAction,
		0,
		dt.OctetString([]byte{aa}),
	)

	desIE := ie.NewInformationElement(
		ie.IEDestinationInterface,
		0,
		dt.OctetString(byte(dInterface)),
	)

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

	createFAR := ie.NewCreateFAR(&farIDIE, &aaIE, &fpIE, nil, nil)
	bCreateFAR, err := createFAR.Serialize()
	if err != nil {
		return nil, err

	}

	createFARIE := ie.NewInformationElement(
		ie.IECreateFAR, //IEcode
		0,              //EntrepriseID
		dt.OctetString(bCreateFAR),
	)
	length = length + ie.IEBasicHeaderSize + createFARIE.Len()
	pfcpHeader := msg.NewPFCPHeader(1, false, true, msg.SessionModificationRequestType, length+12, sei, sn, 0)

	smr := msg.NewPFCPSessionModificationRequest(pfcpHeader, nil, nil, nil, nil, nil, nil, nil, &createPDRIE, &createFARIE, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	return &smr, nil

}
