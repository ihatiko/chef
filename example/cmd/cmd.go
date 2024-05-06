package cmd

import (
	"example/internal/server/deployments/daemon"
	mExample "example/internal/server/deployments/multiple-example"

	"github.com/ihatiko/olymp/core/commands"
)

func Startup() {
	commands.WithApp(
		commands.WithDeployment[daemon.DaemonDeploymentExample](),
		commands.WithDeployment[mExample.MultipleExample](),
	)
}
