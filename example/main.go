package main

import (
	"example/internal/deployments/daemon"
	mExample "example/internal/deployments/multiple-example"
	"github.com/ihatiko/olymp/core/commands"
)

func main() {
	commands.WithApp(
		commands.WithDeployment[daemon.Daemon](),
		commands.WithDeployment[mExample.MultipleExample](),
	)
}
