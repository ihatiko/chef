package peoples

import "example/protoc/peoples"

type ITransport interface {
	peoples.PeoplesServiceServer
}
