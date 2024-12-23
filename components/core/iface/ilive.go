package iface

import (
	"context"
)

type ILive interface {
	Live(ctx context.Context) error
	Name() string
}
