package main

import (
	"example/internal/server/deployments/daemon"
	mExample "example/internal/server/deployments/multiple-example"
	"github.com/ihatiko/olymp/core/commands"
)

func main() {
	commands.WithApp(
		commands.WithDeployment[daemon.Daemon](),
		commands.WithDeployment[mExample.MultipleExample](),
	)
}
