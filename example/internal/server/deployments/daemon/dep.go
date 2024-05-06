package daemon

import (
	planetsReadRepository "example/internal/features/planets/repository/read"
	planetsService "example/internal/features/planets/service"
	planetsTransport "example/internal/features/planets/transport"

	"github.com/ihatiko/olymp/hephaestus/iface"
)

func (d DaemonDeploymentExample) Dep() iface.IDeployment {
	readPostgreSQL := d.ReadPostgreSQL.New()
	readRepository := planetsReadRepository.New(readPostgreSQL)
	service := planetsService.New(readRepository)
	d.transport = planetsTransport.New(service)
	return d
}
