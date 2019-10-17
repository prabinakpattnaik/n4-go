package msg

import (
	"encoding/binary"
)

//The messages inlcuding headers here are implemented as per 3GPPP 29.244 15.5.0 specification

var PFCPBasicHeaderLength = 4

type Message interface {
	Serialize() []byte
	Len() uint16
	Type() PFCPType
}

type PFCPType uint8

const (
	//Node related messages
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

	//Session related messages
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
	SessionEndpointIdentifier uint64
	SequenceNumber            uint32
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
	var b = make([]byte, p.MessageLength+uint16(PFCPBasicHeaderLength))
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
		binary.BigEndian.PutUint64(b[4:12], p.SessionEndpointIdentifier)
		b[12] = uint8(p.SequenceNumber >> 16)
		b[13] = uint8(p.SequenceNumber >> 8)
		b[14] = uint8(p.SequenceNumber)

		b[15] = p.MessagePriority << 4
	} else {
		b[4] = uint8(p.SequenceNumber >> 16)
		b[5] = uint8(p.SequenceNumber >> 8)
		b[6] = uint8(p.SequenceNumber)
		b[7] = 0
	}

	//TODO: mp condition

	return b
}
