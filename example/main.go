package main

import (
	"example/internal/deployments/daemon"
	"example/internal/deployments/grpc"
	"github.com/ihatiko/chef/core/commands"
)

func main() {
	commands.WithApp(
		commands.WithDeployment[daemon.Deployment](),
		commands.WithDeployment[grpc.Deployment](),
	)
}
