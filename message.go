package n4

//29.244 15.5.0

type PFCPType uint8

const (
	//Node related messages
	HeartbeatRequest            PFCPType = 1
	HeartbeatResponse           PFCPType = 2
	PFDManagementRequest        PFCPType = 3
	PFDManagementResponse       PFCPType = 4
	AssociationSetupRequest     PFCPType = 5
	AssociationSetupResponse    PFCPType = 6
	AssociationUpdateRequest    PFCPType = 7
	AssociationUpdateResponse   PFCPType = 8
	AssociationReleaseRequest   PFCPType = 9
	AssociationReleaseResponse  PFCPType = 10
	VersionNotSupportedResponse PFCPType = 11
	NodeReportRequest           PFCPType = 12
	NodeReportResponse          PFCPType = 13

	//Session related messages
	SessionEstablishmentRequest  PFCPType = 50
	SessionEstablishmentResponse PFCPType = 51
	SessionModificationRequest   PFCPType = 52
	SessionModificationResponse  PFCPType = 53
	SessionDeletionRequest       PFCPType = 54
	SessionDeletionResponse      PFCPType = 55
	SessionReportRequest         PFCPType = 56
	SessionReportResponse        PFCPType = 57
)

// Message represents the COAP message
type PFCPHeader struct {
	Version                  uint8
	MP                       bool
	S                        bool
	MessageType              PFCPType
	MessageLength            uint16
	SequenceNumber           uint32
	SessionEndpointIdentider uint64
	messagePriority          uint8
}

type PFCPMessage struct {
	Header PFCPHeader
}
