package msg

import (
	"encoding/binary"
	"fmt"

	"github.com/prabinakpattnaik/n4-go/ie"
	"github.com/u-root/u-root/pkg/uio"
)

//The messages inlcuding headers here are implemented as per 3GPPP 29.244 15.5.0 specification

const ( PFCPVERSION   uint8 = 1

	SEID_PRESENT        = true
	MP_PRESENT          = true
)


var (
	PFCPBasicHeaderLength = 4
	PFCPBasicMessageSize  = uint16(4)
	PFCPMessageSize       = uint16(12)
)

type PFCPType uint8

const (
	HeartbeatRequestType            PFCPType = 1
	HeartbeatResponseType           PFCPType = 2
	PFDManagementRequestType        PFCPType = 3
	PFDManagementResponseType       PFCPType = 4
	AssociationSetupRequestType     PFCPType = 5
	AssociationSetupResponseType    PFCPType = 6
	AssociationUpdateRequestType    PFCPType = 7
	AssociationUpdateResponseType   PFCPType = 8
	AssociationReleaseRequestType   PFCPType = 9
	AssociationReleaseResponseType  PFCPType = 10
	VersionNotSupportedResponseType PFCPType = 11
	NodeReportRequestType           PFCPType = 12
	NodeReportResponseType          PFCPType = 13

	SessionEstablishmentRequestType  PFCPType = 50
	SessionEstablishmentResponseType PFCPType = 51
	SessionModificationRequestType   PFCPType = 52
	SessionModificationResponseType  PFCPType = 53
	SessionDeletionRequestType       PFCPType = 54
	SessionDeletionResponseType      PFCPType = 55
	SessionReportRequestType         PFCPType = 56
	SessionReportResponseType        PFCPType = 57
)

// Message represents the COAP message
type PFCPHeader struct {
	Version                   uint8
	MP                        bool
	S                         bool
	MessageType               PFCPType
	MessageLength             uint16
	SequenceNumber            uint32
	SessionEndpointIdentifier uint64
	MessagePriority           uint8
}

func NewPFCPHeader(v uint8, mp bool, s bool, mt PFCPType, ml uint16, sei uint64, sn uint32, messagepriority uint8) *PFCPHeader {

	return &PFCPHeader{
		Version:                   v,
		MP:                        mp,
		S:                         s,
		MessageType:               mt,
		MessageLength:             ml,
		SessionEndpointIdentifier: sei,
		SequenceNumber:            sn,
		MessagePriority:           messagepriority,
	}

}

func (p PFCPHeader) Serialize() []byte {
	var b = make([]byte, uint16(PFCPBasicHeaderLength))
	var a uint8

	//first byte format: [v][v][v][][][][mp][s]
	if p.S {
		a = 1
	}

	if p.MP {
		a |= 2
	}

	a |= (p.Version << 5)
	b[0] = a

	b[1] = byte(p.MessageType)
	binary.BigEndian.PutUint16(b[2:4], p.MessageLength)

	if p.S {
		buf :=make([]byte,8)	
		binary.BigEndian.PutUint64(buf, p.SessionEndpointIdentifier)
		b=append(b,buf...)
		b=append(b, uint8(p.SequenceNumber >> 16))
		b=append(b, uint8(p.SequenceNumber >> 8))
		b=append(b, uint8(p.SequenceNumber))
		b=append(b, p.MessagePriority << 4)
	} else {
		b=append(b, uint8(p.SequenceNumber >> 16))
		b=append(b,uint8(p.SequenceNumber >> 8))
		b=append(b, uint8(p.SequenceNumber))
		b=append(b, 0)
	}

	//TODO: mp condition

	return b

}

type PFCPMessage struct {
	Header *PFCPHeader
	IEs    ie.InformationElements
}

//TODO:Basic interface
type PFCP interface {
	Type() PFCPType
	Serialize() ([]byte, error)
	GetHeader() *PFCPHeader
	//GetIEs() []ie.InformationElement
	//AddIE(i ie.InformationElement)
}

// MessageFromBytes parses a PFCPMessage message from a byte stream.
func MessageFromBytes(data []byte) (*PFCPMessage, error) {
	buf := uio.NewBigEndianBuffer(data)
	f := buf.Read8()
	s := 0x01 & f
	mp := 0x02 & f
	v := uint8(f >> 5)

	pfcpMessageType := PFCPType(buf.Read8())

	if !((pfcpMessageType >= HeartbeatRequestType && pfcpMessageType <= NodeReportResponseType) || (pfcpMessageType >= SessionEstablishmentRequestType && pfcpMessageType <= SessionReportResponseType)) {
		return nil, fmt.Errorf("wrong message type")
	}

	m, err := buf.ReadN(2)
	if err != nil {
		return nil, fmt.Errorf("wrong message length")
	}

	ml := binary.BigEndian.Uint16(m)
	var sei uint64
	if s == 1 {

		seiBytes, err := buf.ReadN(8)
		if err != nil {
			return nil, fmt.Errorf("wrong session endpoint identifer")
		}
		sei = binary.BigEndian.Uint64(seiBytes)

	}

	snByte, err := buf.ReadN(3)
	if err != nil {
		return nil, fmt.Errorf("wrong session number")
	}

	snBytes := []byte{0x00}
	var sn uint32
	snBytes = append(snBytes, snByte...)
	sn = binary.BigEndian.Uint32(snBytes)

	var messagepriority uint8
	if mp == 2 {
		messagepriority = buf.Read8() | 0xF0
	} else {
		_ = buf.Read8()
	}

	pfcpHeader := NewPFCPHeader(v, (mp == 2), (s == 1), pfcpMessageType, ml, sei, sn, messagepriority)

	p := &PFCPMessage{
		Header: pfcpHeader,
	}

	if err := p.IEs.FromBytes(buf.Data()); err != nil {
		return nil, err
	}
	return p, nil
}
