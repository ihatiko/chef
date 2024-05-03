package iface

type IPkg[T any] interface {
	New() (T, error)
}
