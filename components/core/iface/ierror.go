package iface

type IPkg interface {
	Error() error
	Name() string
	HasError() bool
}
