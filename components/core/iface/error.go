package iface

type IError interface {
	HasError() bool
	Error() error
	String() string
}
