package daemon

import (
	"github.com/ihatiko/olymp/hephaestus/iface"
)

func (d DaemonDeploymentExample) Dep() iface.IDeployment {
	//readPostgreSQL := d.ReadPostgreSQL.New()
	//readRepository := planetsReadRepository.New(readPostgreSQL)
	//service := planetsService.New(readRepository)
	//d.transport = planetsTransport.New(service)
	return d
}
