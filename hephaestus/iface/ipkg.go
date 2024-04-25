package iface

import "context"

type IPkg[T any] interface {
	New() (T, error)
	Live(context.Context, T) bool
}
