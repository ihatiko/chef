package main

import (
	daemonDeployment "example/internal/deployments/daemon"
	grpcDeployment "example/internal/deployments/grpc"
	multipleExampleDeployment "example/internal/deployments/multiple_example"
	"github.com/ihatiko/olymp/core/commands"
)

func main() {
	commands.WithApp(
		commands.WithDeployment[daemonDeployment.Daemon](),
		commands.WithDeployment[multipleExampleDeployment.MultipleExample](),
		commands.WithDeployment[grpcDeployment.GrpcExample](),
	)
}
