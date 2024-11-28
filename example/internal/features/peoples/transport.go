package peoples

import (
	"example/pkg/protoc/peoples"
)

type ITransport interface {
	peoples.PeoplesServiceServer
}
