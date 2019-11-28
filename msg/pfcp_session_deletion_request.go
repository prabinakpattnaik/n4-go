package msg

import (
	"fmt"
)

//PFCPSessionDeletionRequest
type PFCPSessionDeletionRequest struct {
	Header *PFCPHeader
}

//NewPFCPSessionDeletionRequest creates new PFCPSessionDeletionRequst
func NewPFCPSessionDeletionRequest(h *PFCPHeader) PFCPSessionDeletionRequest {
	return PFCPSessionDeletionRequest{
		Header: h,
	}
}
func (sdr PFCPSessionDeletionRequest) Serialize() ([]byte, error) {
	fmt.Printf("[%x]\n", sdr.Header.Serialize())
	return sdr.Header.Serialize(), nil
}

func (sdr PFCPSessionDeletionRequest) Len() uint16 {
	return uint16(PFCPBasicHeaderLength) + sdr.Header.MessageLength
}

func (sdr PFCPSessionDeletionRequest) Type() PFCPType {
	return sdr.Header.MessageType
}

func (sdr PFCPSessionDeletionRequest) GetHeader() *PFCPHeader {
	return sdr.Header
}
