package store

import (
	"log/slog"
	"sync"

	"github.com/ihatiko/chef/components/core/iface"
)

var PackageStore = errorStore{
	mt:       sync.Mutex{},
	packages: []iface.IPkg{},
}

type errorStore struct {
	mt       sync.Mutex
	packages []iface.IPkg
}

func (s *errorStore) Load(e iface.IPkg) {
	s.mt.Lock()
	defer s.mt.Unlock()
	if !e.HasError() {
		slog.Info("successfully package connect", slog.String("component", e.Name()))
	}
	s.packages = append(s.packages, e)
}

func (s *errorStore) Get() []iface.IPkg {
	return s.packages
}
