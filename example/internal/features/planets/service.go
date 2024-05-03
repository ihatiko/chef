package planets

import (
	"github.com/ihatiko/olymp/temple/transports/cron"
	"github.com/ihatiko/olymp/temple/transports/daemon"
)

type IService interface {
	Load(request daemon.Request) error
	Update(request cron.Request) error
}
