package ie

type CauseID uint8


const (

RequestAccepted                 CauseID =1
RequestRejected                 CauseID =64
SessionContextNotFound          CauseID =65
MandatoryIEMissing              CauseID =66
ConditionalIEMissing            CauseID =67
InvalidLength                   CauseID =68
MandatoryIEIncorrect            CauseID =69
InvalidForwardingPolicy         CauseID =70
InvalidFTEIDAllocationOption    CauseID =71
NoEstablishedPFPCPAssociation   CauseID =72
RuleCreationModificationFailure CauseID =73
PFCPEntityInCongestion          CauseID =74
NoResourcesAvailable            CauseID =75
ServiceNotSupported             CauseID =76
SystemFailure                   CauseID =77

)
