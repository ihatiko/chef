package iface

import "context"

type IPkg[T any] interface {
	NewClient() (T, error)
	Live(context.Context, T) bool
}
