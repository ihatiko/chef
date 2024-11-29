//go:generate deployment-dependency

package multiple_example

import (
	peoplesRepository "example/internal/features/peoples/repository"
	peoplesService "example/internal/features/peoples/service"
	peoplesTransport "example/internal/features/peoples/transport"
	planetsReadRepository "example/internal/features/planets/repository/read"
	planetsService "example/internal/features/planets/service"
	planetsTransport "example/internal/features/planets/transport"

	"github.com/ihatiko/olymp/core/iface"
)

func (d MultipleExample) Dep() iface.IDeployment {
	readPostgreSQL := d.ReadPostgreSQL.New()
	planetsRR := planetsReadRepository.New(readPostgreSQL)
	planetsS := planetsService.New(planetsRR)
	d.iPlanetsTransport = planetsTransport.New(planetsS)

	redis := d.Redis.New()
	peoplesR := peoplesRepository.New(redis)
	peoplesS := peoplesService.New(peoplesR)
	d.iPeoplesTransport = peoplesTransport.New(peoplesS)
	return d
}