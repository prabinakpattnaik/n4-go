package msg

import (
	"fmt"

	"github.com/prabinakpattnaik/n4-go/ie"
	dt "github.com/fiorix/go-diameter/diam/datatype"
	log "github.com/sirupsen/logrus"
)

//PFCPSessionModificationRequest
type PFCPSessionModificationRequest struct {
	Header                   *PFCPHeader
	CPFSEID                  *ie.InformationElement
	RemovePDR                *ie.InformationElement
	RemoveFAR                *ie.InformationElement
	RemoveURR                *ie.InformationElement
	RemoveQER                *ie.InformationElement
	RemoveBAR                *ie.InformationElement
	RemoveTrafficEndpoint    *ie.InformationElement
	CreatePDR                *ie.InformationElements
	CreateFAR                *ie.InformationElements
	CreateURR                *ie.InformationElement
	CreateQER                *ie.InformationElement
	CreateBAR                *ie.InformationElement
	CreateTrafficEndpoint    *ie.InformationElement
	UpdatePDR                *ie.InformationElement
	UpdateFAR                *ie.InformationElement
	UpdateURR                *ie.InformationElement
	UpdateQER                *ie.InformationElement
	UpdateBAR                *ie.InformationElement
	UpdateTrafficEndpoint    *ie.InformationElement
	PFCPSMReqFlags           *ie.InformationElement
	QueryURR                 *ie.InformationElement
	UserPlaneInactivityTimer *ie.InformationElement
	QueryURRReference        *ie.InformationElement
	TraceInformation         *ie.InformationElement
}

//NewPFCPSessionModificationRequest creates new PFCPSessionMondificationRequst
func NewPFCPSessionModificationRequest(h *PFCPHeader, cpfseid, rPDR, rFAR, rURR, rQER, rBAR, rTE *ie.InformationElement, cPDR, cFAR *ie.InformationElements, cURR, cQER, cBAR, cTE, uPDR, uFAR, uURR, uQER, uBAR, uTE, pfcpSMReqFlags, qURR, userPIT, qURRReference, traceInformation *ie.InformationElement) PFCPSessionModificationRequest {
	return PFCPSessionModificationRequest{
		Header:                   h,
		CPFSEID:                  cpfseid,
		RemovePDR:                rPDR,
		RemoveFAR:                rFAR,
		RemoveURR:                rURR,
		RemoveQER:                rQER,
		RemoveBAR:                rBAR,
		RemoveTrafficEndpoint:    rTE,
		CreatePDR:                cPDR,
		CreateFAR:                cFAR,
		CreateURR:                cURR,
		CreateQER:                cQER,
		CreateBAR:                cBAR,
		CreateTrafficEndpoint:    cTE,
		UpdatePDR:                uPDR,
		UpdateFAR:                uFAR,
		UpdateURR:                uURR,
		UpdateQER:                uQER,
		UpdateBAR:                uBAR,
		UpdateTrafficEndpoint:    uTE,
		PFCPSMReqFlags:           pfcpSMReqFlags,
		QueryURR:                 qURR,
		UserPlaneInactivityTimer: userPIT,
		QueryURRReference:        qURRReference,
		TraceInformation:         traceInformation,
	}

}

func (smr PFCPSessionModificationRequest) Serialize() ([]byte, error) {
	messageLength := smr.GetHeader().MessageLength+ uint16(PFCPBasicHeaderLength) 
	output := smr.Header.Serialize()
	b := []byte{}
	var err error
	if smr.CPFSEID != nil {
		b, err = smr.CPFSEID.Serialize()
		if err != nil {
			return nil, err
		}

	}

	if smr.RemovePDR != nil {
		rpdr, err := smr.RemovePDR.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, rpdr...)
	}

	if smr.RemoveFAR != nil {
		rfar, err := smr.RemoveFAR.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, rfar...)
	}

	if smr.RemoveURR != nil {
		rurr, err := smr.RemoveURR.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, rurr...)
	}

	if smr.RemoveQER != nil {
		rqer, err := smr.RemoveQER.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, rqer...)
	}

	if smr.RemoveBAR != nil {
		rbar, err := smr.RemoveBAR.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, rbar...)
	}

	if smr.RemoveTrafficEndpoint != nil {
		rte, err := smr.RemoveTrafficEndpoint.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, rte...)
	}
        /*
	if smr.CreatePDR != nil {
		cpdr, err := smr.CreatePDR.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, cpdr...)
	}

	if smr.CreateFAR != nil {
		cfar, err := smr.CreateFAR.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, cfar...)
	}
	
	for _, pdrie := range *smr.CreatePDR {
                createpdr, err := pdrie.Serialize()
                if err != nil {
                        return nil, err
                }
               b = append(b, createpdr...)
        }

        for _, farie := range *smr.CreateFAR {
                 createfar, err := farie.Serialize()
                if err != nil {
                        
                        return nil, err
                }
                b = append(b, createfar...)
        }
*/

	if smr.CreateURR != nil {
		curr, err := smr.CreateURR.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, curr...)
	}

	if smr.CreateQER != nil {
		cqer, err := smr.CreateQER.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, cqer...)
	}

	if smr.CreateBAR != nil {
		cbar, err := smr.CreateBAR.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, cbar...)
	}

	if smr.CreateTrafficEndpoint != nil {
		cte, err := smr.CreateTrafficEndpoint.Serialize()
		if err != nil {
			return nil, err
		}
		b = append(b, cte...)
	}
	if smr.UpdatePDR != nil {
		log.WithFields(log.Fields{"data": smr.UpdatePDR}).Info("received UpdatePDR")
                updr , err := smr.UpdatePDR.Serialize()
		log.WithFields(log.Fields{"data": updr}).Info("serialization UpdatePDR")
                if err != nil {
			log.WithError(err).Error("error in UpdatePDR serialization")
                        return nil, err
                }
                b = append(b, updr...)
        }

        if smr.UpdateFAR != nil {
                ufar, err := smr.UpdateFAR.Serialize()
                if err != nil {
			log.WithError(err).Error("error in UpdateFAR serialization")
                        return nil, err
                }
                b = append(b, ufar...)
        }

        if smr.UpdateURR != nil {
                uurr, err := smr.UpdateURR.Serialize()
                if err != nil {
                        return nil, err
                }
                b = append(b, uurr...)
        }
	if smr.UpdateQER != nil {
                uqer, err := smr.UpdateQER.Serialize()
                if err != nil {
                        return nil, err
                }
                b = append(b, uqer...)
        }

        if smr.UpdateBAR != nil {
                ubar, err := smr.UpdateBAR.Serialize()
                if err != nil {
                        return nil, err
                }
                b = append(b, ubar...)
        }

	output=append(output, b...)
	if uint16(len(output))==messageLength {
		return output, nil
	}else {
	//fmt.Printf("left length is %u\n", uint16(len(output)))
	//fmt.Printf("messageLength length is %u\n", messageLength)
	return nil, fmt.Errorf("Error in serialization of PFCP Session Modification Request")
}

}

//ProcessPFCPSessionModificationRequest process session modification request and produce the result (either SessionModificationResponse or error)
func ProcessPFCPSessionModificationRequest(m *PFCPMessage, sSEID uint64) ([]byte, error) {
	//check if F-SEID is unknown
	//rule is not there
	//sucess
	pfcpMessage, err := FromPFCPMessage(m)
	if err != nil {
		return nil, err
	}
	pfcpSessionModificationRequest, ok := pfcpMessage.(PFCPSessionModificationRequest)

	if !ok {
		return nil, fmt.Errorf("PFCP Session EstablishmentModification Request could not type asserted")
	}

	c := ie.NewInformationElement(
		ie.IECause,
		0,
		dt.OctetString([]byte{0x01}),
	)
	length := ie.IEBasicHeaderSize + c.Len()
	//how to select SEID from sending side.
	pfcpHeader := NewPFCPHeader(1, false, true, SessionModificationResponseType, length, sSEID, pfcpSessionModificationRequest.Header.SequenceNumber, 0)
	pfcpSessionModificationResponse := NewPFCPSessionModificationResponse(pfcpHeader, &c, nil, nil)
	b, err := pfcpSessionModificationResponse.Serialize()
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (smr PFCPSessionModificationRequest) Len() uint16 {
	return uint16(PFCPBasicHeaderLength) + smr.Header.MessageLength
}

func (smr PFCPSessionModificationRequest) Type() PFCPType {
	return smr.Header.MessageType
}

func (smr PFCPSessionModificationRequest) GetHeader() *PFCPHeader {
	return smr.Header
}
