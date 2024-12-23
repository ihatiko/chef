package planets

import (
	"github.com/ihatiko/chef/components/transports/cron"
	"github.com/ihatiko/chef/components/transports/daemon"
)

type IService interface {
	Load(request daemon.Request) error
	Update(request cron.Request) error
}
