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
	planetsReadRepository := planetsReadRepository.New(readPostgreSQL)
	planetsService := planetsService.New(planetsReadRepository)
	d.iPlanetsTransport = planetsTransport.New(planetsService)

	redis := d.Redis.New()
	peoplesRepository := peoplesRepository.New(redis)
	peoplesService := peoplesService.New(peoplesRepository)
	d.iPeoplesTransport = peoplesTransport.New(peoplesService)
	return d
}
