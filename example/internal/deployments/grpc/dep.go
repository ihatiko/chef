//go:generate deployment-dependency

package grpc

import (
	"github.com/ihatiko/olymp/core/iface"
)

func (d GrpcExample) Dep() iface.IDeployment {
	//readPostgreSQL := d.ReadPostgreSQL.New()
	//readRepository := planetsReadRepository.New(readPostgreSQL)
	//service := planetsService.New(readRepository)
	//d.transport = planetsTransport.New(service)
	return d
}
