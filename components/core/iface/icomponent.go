package iface

type IComponent interface {
	IShutdownComponent
	ILive
	Run() error
}
