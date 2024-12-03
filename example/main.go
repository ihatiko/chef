package main

import (
	daemonDeployment "example/internal/deployments/daemon"
	grpcDeployment "example/internal/deployments/grpc"
	"github.com/ihatiko/olymp/core/commands"
)

func main() {
	commands.WithApp(
		commands.WithDeployment[daemonDeployment.Daemon](),
		commands.WithDeployment[grpcDeployment.Grpc](),
	)
}
