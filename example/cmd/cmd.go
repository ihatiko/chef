package cmd

import (
	"example/internal/server/deployments/daemon"

	"github.com/ihatiko/olymp/hephaestus/commands"
)

func Startup() {
	commands.WithApp(
		commands.WithDeployment[daemon.DaemonDeploymentExample](),
	)
}
