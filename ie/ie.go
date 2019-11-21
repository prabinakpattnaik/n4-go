package ie

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/u-root/u-root/pkg/uio"
)

type IEType uint16

const (
	EnterpriseIDStart = 32768
	IEBasicHeaderSize = 4
)

const (

	//3GPP TS 29.244 Table 8.1.2-1: Information Element Types
	IEReserved                        IEType = 0
	IECreatePDR                       IEType = 1
	IEPDI                             IEType = 2
	IECreateFAR                       IEType = 3
	IEForwardingParameters            IEType = 4
	IEDuplicatingParameters           IEType = 5
	IECreateURR                       IEType = 6
	IECreateQER                       IEType = 7
	IECreatedPDR                      IEType = 8
	IEUpdatePDR                       IEType = 9
	IEUpdateFAR                       IEType = 10
	IEUpdateForwardingParamets        IEType = 11
	IEUpdateBAR                       IEType = 12 //PFCP Session Report Response
	IEUpdateURR                       IEType = 13
	IEUpdateQER                       IEType = 14
	IERemovePDR                       IEType = 15
	IERemoveFAR                       IEType = 16
	IERemoveURR                       IEType = 17
	IERemoveQER                       IEType = 18
	IECause                           IEType = 19
	IESourceInterface                 IEType = 20
	IEFTEID                           IEType = 21
	IENetworkInstance                 IEType = 22
	IESDFFilter                       IEType = 23
	IEApplicationID                   IEType = 24
	IEGateStatus                      IEType = 25
	IEMBR                             IEType = 26
	IEGBR                             IEType = 27
	IEQERCorrelationID                IEType = 28
	IEPrecedence                      IEType = 29
	IETransportLevelMarking           IEType = 30
	IEVolumeThreshold                 IEType = 31
	IETimeThreshold                   IEType = 32
	IEMonitoringTime                  IEType = 33
	IESubsequentVolumeThreshold       IEType = 34
	IESubsequentTimeThreshold         IEType = 35
	IEInactivityDetectionTime         IEType = 36
	IEReportingTriggers               IEType = 37
	IERedirectInformation             IEType = 38
	IEReportType                      IEType = 39
	IEOffendingIE                     IEType = 40
	IEForwardingPolicy                IEType = 41
	IEDestinationInterface            IEType = 42
	IEUPFunctionFeatures              IEType = 43
	IEApplyAction                     IEType = 44
	IEDownlinkDataServiceInformation  IEType = 45
	IEDownlinkDataNotificationDelay   IEType = 46
	IEDLBufferingDuration             IEType = 47
	IEDLBufferingSuggestedPacketCount IEType = 48
	IEPFCPSMReqFlags                  IEType = 49
	IEPFCPSRRspFlags                  IEType = 50

	IELoadControlInformation          IEType = 51
	IESequenceNumber                  IEType = 52
	IEMetric                          IEType = 53
	IEOverloadControlInformation      IEType = 54
	IETimer                           IEType = 55
	IEPDRID                           IEType = 56
	IEFSEID                           IEType = 57
	IEApplicationIDPFDs               IEType = 58
	IEPFDContext                      IEType = 59
	IENodeID                          IEType = 60
	IEPFDContents                     IEType = 61
	IEMeasurementMethod               IEType = 62
	IEUsageReportTrigger              IEType = 63
	IEMeasurementPeriod               IEType = 64
	IEFQCSID                          IEType = 65
	IEVolumeMeasurement               IEType = 66
	IEDurationMeasurement             IEType = 67
	IEApplicationDetectionInformation IEType = 68
	IETimeOfFirstPacket               IEType = 69
	IETimeofLastPacket                IEType = 70
	IEQuotaHoldingTime                IEType = 71
	IEDroppedDLTrafficThreshold       IEType = 72
	IEVolumeQuota                     IEType = 73
	IETimeQuota                       IEType = 74
	IEStartTime                       IEType = 75
	IEEndTime                         IEType = 76
	IEQueryURR                        IEType = 77
	IEUsageReportSMR                  IEType = 78
	IEUsageReportSDR                  IEType = 79
	IEUsageReportSRR                  IEType = 80
	IEURRID                           IEType = 81
	IELinkedURRID                     IEType = 82
	IEDownlinkDataReport              IEType = 83
	IEOuterHeaderCreation             IEType = 84
	IECreateBAR                       IEType = 85
	IEUpdateBARSMR                    IEType = 86
	IERemoveBAR                       IEType = 87
	IEBARID                           IEType = 88
	IECPFunctionFeatures              IEType = 89
	IEUsageInformation                IEType = 90
	IEApplicationInstanceID           IEType = 91
	IEFlowInformation                 IEType = 92
	IEUEIPaddress                     IEType = 93
	IEPacketRate                      IEType = 94
	IEOuterHeaderRemoval              IEType = 95
	IERecoveryTimestamp               IEType = 96
	IEDLFLowLevelMarking              IEType = 97
	IEHeaderEnrichment                IEType = 98

	IEMeasurementInformation  IEType = 100
	IENodeReportType          IEType = 101
	IERemoteGTPUPeer          IEType = 103
	IEUESEQN                  IEType = 104
	IEActivatePredefinedRules IEType = 106
	IEDeactivePredeinedRules  IEType = 107
	IEFARID                   IEType = 108
	IEQERID                   IEType = 109

	IEUserPlaneIPResourceInformation IEType = 116
)

//Informatin Element is a Type, length,value group
//bit 8 of octet 1 is not set, 3GPP defined by and Entreprise ID is absent. otherwise it has Entreprise ID
type InformationElement struct {
	Type         IEType
	Length       uint16
	EnterpriseID uint16
	Data         datatype.Type
}

//NewInformation element creates new IE

func NewInformationElement(ietype IEType, enterpriseId uint16, data datatype.Type) InformationElement {
	i := InformationElement{
		Type:         ietype,
		EnterpriseID: enterpriseId,
		Data:         data,
	}
	i.Length = i.Len()
	return i

}

//
func DecodeIE(data []byte) (*InformationElement, error) {
	i := &InformationElement{}
	if err := i.DecodeFromBytes(data); err != nil {
		return i, err
	}
	return i, nil
}

//DecodeFromBytes decodes the bytes of Information element
func (i *InformationElement) DecodeFromBytes(data []byte) error {

	if len(data) < 4 {
		return fmt.Errorf("Not enough data to decode InformationElement header: %d bytes", len(data))
	}

	i.Type = IEType(binary.BigEndian.Uint16(data[0:2]))
	i.Length = binary.BigEndian.Uint16(data[3:4])

	if i.Type >= EnterpriseIDStart {
		i.EnterpriseID = binary.BigEndian.Uint16(data[5:6])
	}
	var err error
	i.Data, err = datatype.DecodeOctetString(data[:i.Length])

	return err
}

//Serialize returns of byte sequence of this Information Element
func (i *InformationElement) Serialize() ([]byte, error) {
	if i.Data == nil {
		return nil, errors.New("Failed to serialize Information Element: Data is nil")
	}

	b := make([]byte, IEBasicHeaderSize+i.Len())
	err := i.SerializeTo(b)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (i *InformationElement) SerializeTo(b []byte) error {
	binary.BigEndian.PutUint16(b[0:2], uint16(i.Type))
	hl := i.Len()
	copy(b[2:4], []byte{uint8(hl >> 8), uint8(hl)})
	payload := i.Data.Serialize()
	if i.Type < EnterpriseIDStart {
		copy(b[4:], payload)
	}

	// reset padding bytes
	/*
		b = b[hl+len(payload):]
		for i := 0; i < a.Data.Padding(); i++ {
			b[i] = 0
		}
	*/
	return nil
}

//Len returns the length of Inforamtion Element
func (i *InformationElement) Len() uint16 {
	//TODO include entreprise ID
	return uint16(i.Data.Len())
}

func (i *InformationElement) String() string {
	return fmt.Sprintf("{Code:%d,Length:%d,EntrepriseID:%d,Data:%s}", i.Type, i.Length, i.EnterpriseID, i.Data)
}

type InformationElements []InformationElement

// FromBytes reads data into ies and returns an error if the ies are not a
// valid serialized representation of PFCP Information Elements
func (ies *InformationElements) FromBytes(data []byte) error {
	//TODO: 10 is not good value.
	*ies = make(InformationElements, 0, 10)
	if len(data) == 0 {
		// no InformationElements
		return nil
	}
	buf := uio.NewBigEndianBuffer(data)
	var ie InformationElement
	for buf.Has(4) {
		ietype := IEType(buf.Read16())
		length := buf.Read16()
		data, err := buf.ReadN(int(length))
		if err != nil {
			return err
		}

		switch ietype {
		case IERecoveryTimestamp:
			v, err := datatype.DecodeTime(data)
			if err != nil {
				return err
			}
			ie = NewInformationElement(ietype, 0, v)
		default:
			// TODO: parse the data into specific type.
			ie = NewInformationElement(ietype, 0, datatype.OctetString(data))
		}

		*ies = append(*ies, ie)
	}
	return buf.FinError()

}

//
