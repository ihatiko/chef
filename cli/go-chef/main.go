package main

import (
	commandExecutor "github.com/ihatiko/chef/components/code-gen-utils/command-executor"
	packageUpdate "github.com/ihatiko/chef/components/code-gen-utils/package-updater"
	"log/slog"
	"os"
	"strings"
)

const corePathPackage = "github.com/ihatiko/gochef/cli/go-chef-core"
const coreNamePackage = "go-chef-core"

func main() {
	params := strings.Join(os.Args[1:], " ")
	//TODO timeout on update
	packageUpdate.AutoUpdate(corePathPackage)
	composer := commandExecutor.NewExecutor()
	result, err := composer.ExecDefaultCommand(coreNamePackage)
	if err != nil {
		slog.Error("Error executing command: ", slog.Any("error", err), slog.String("command", params))
	}
	slog.Info(result.String())
}
