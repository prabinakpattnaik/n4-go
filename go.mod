module github.com/prabinakpattnaik/n4-go

go 1.16

replace (
	magma/feg/cloud/go => ../../../feg/cloud/go
	magma/feg/cloud/go/protos => ../../../feg/cloud/go/protos
	magma/gateway => ../../../orc8r/gateway/go
	magma/lte/cloud/go => ../../../lte/cloud/go
	magma/orc8r/cloud/go => ../../../orc8r/cloud/go
	magma/orc8r/lib/go => ../../../orc8r/lib/go
	magma/orc8r/lib/go/protos => ../../../orc8r/lib/go/protos
)

require (
	github.com/fiorix/go-diameter v3.0.2+incompatible
	github.com/golang/glog v1.1.1
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/u-root/u-root v6.0.0+incompatible
	github.com/urfave/cli v1.22.3
	google.golang.org/genproto v0.0.0-20230403163135-c38d8f061ccd // indirect
	google.golang.org/grpc v1.54.0
	magma/feg/cloud/go v0.0.0
	magma/feg/cloud/go/protos v0.0.0
	magma/lte/cloud/go v0.0.0
	magma/orc8r/cloud/go v0.0.0
	magma/orc8r/lib/go v0.0.0 // indirect
	magma/orc8r/lib/go/protos v0.0.0

)
