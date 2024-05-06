package iface

type IDeployment interface {
	Run()
	Dep() IDeployment
}
