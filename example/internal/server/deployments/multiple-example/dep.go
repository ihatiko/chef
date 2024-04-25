package daemon

import "github.com/ihatiko/olymp/hephaestus/iface"

func (d MultipleExample) Dep() iface.IDeployment {
	return d
}
