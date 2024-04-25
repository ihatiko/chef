package iface

import (
	"context"
	"time"
)

type IBeforeLifecycleComponent interface {
	BeforeShutdown() error
	Name() string
}

type IAfterLifecycleComponent interface {
	AfterShutdown() error
	Name() string
}

type IShutdownComponent interface {
	TimeToWait() time.Duration
	Shutdown() error
	Name() string
}

type ILive interface {
	Live(ctx context.Context) error
	Name() string
}
