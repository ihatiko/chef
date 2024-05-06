package planets

import (
	"github.com/ihatiko/olymp/components/transports/cron"
	"github.com/ihatiko/olymp/components/transports/daemon"
)

type IService interface {
	Load(request daemon.Request) error
	Update(request cron.Request) error
}
