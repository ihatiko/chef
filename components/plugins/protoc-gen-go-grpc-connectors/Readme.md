```bash
go install github.com/ihatiko/chef/protoc-gen-go-connectors

```


### Result
```go
// GENERATED BY PROTOC-GEN-GO-CONNECTORS
// DOES NOT EDIT

package planets

import "github.com/ihatiko/chef/components/transports/grpc"

const (
	sdkGrpcName = "GrpcExampleDeployment PlanetsService"
)

type PlanetsConfig struct {
	grpc.Config
}

type PlanetsTransport struct {
	grpc.Transport
}

func (p *PlanetsConfig) Use() PlanetsTransport {
	return PlanetsTransport{Transport: p.Config.Use()}
}

func (p *PlanetsConfig) Name() string {
	return sdkGrpcName
}

func (t PlanetsTransport) Routing(impl PlanetsServiceServer) PlanetsTransport {
	t.Transport.Routing(PlanetsService_ServiceDesc, impl)
	return t
}

```