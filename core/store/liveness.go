package store

import (
	"sync"

	"github.com/ihatiko/olymp/core/iface"
)

var LivenessStore = store{
	mt:         sync.Mutex{},
	components: []iface.ILive{},
}

type store struct {
	mt         sync.Mutex
	components []iface.ILive
}

func (s store) Load(live iface.ILive) {
	s.mt.Lock()
	defer s.mt.Unlock()
	s.components = append(s.components, live)
}

func (s store) Get() []iface.ILive {
	return s.components
}
