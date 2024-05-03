package monads

type IError interface {
	HasError() bool
	Error() error
	String() string
}
