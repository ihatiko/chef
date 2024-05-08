package kafka

import (
	"strings"
	"sync"
)

var transport map[string]Transport = make(map[string]Transport)
var active map[string]struct{} = make(map[string]struct{})
var shutdown map[string]struct{} = make(map[string]struct{})
var mt sync.Mutex

type Transport struct {
	cfg Config
}

func (cfg Config) Use() Transport {
	mt.Lock()
	key := strings.Join(cfg.Brokers, ",")
	if t, ok := transport[key]; ok {
		defer mt.Unlock()
		return t
	}
	t := new(Transport)
	t.cfg = cfg
	transport[key] = *t
	return *t
}

func (t Transport) Run() {

}
