package ie

import (
	"github.com/fiorix/go-diameter/diam/datatype"
)

var CodeDataTypeMapping = map[IEType]datatype.TypeID{
	96: datatype.OctetStringType,
}
