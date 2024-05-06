package store

import (
	"sync"

	"github.com/ihatiko/olymp/core/iface"
)

var PackageStore = errorStore{
	mt:     sync.Mutex{},
	errors: []iface.IPkg{},
}

type errorStore struct {
	mt     sync.Mutex
	errors []iface.IPkg
}

func (s *errorStore) Load(e iface.IPkg) {
	s.mt.Lock()
	defer s.mt.Unlock()
	s.errors = append(s.errors, e)
}

func (s *errorStore) Get() []iface.IPkg {
	return s.errors
}
