//go:generate deployment-dependency

package multiple_example

import "github.com/ihatiko/olymp/core/iface"

func (d MultipleExample) Dep() iface.IDeployment {

	return d
}
