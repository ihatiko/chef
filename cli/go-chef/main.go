package main

import (
	filebuilder "github.com/ihatiko/chef/code-gen/file-builder"
	packageUpdate "github.com/ihatiko/chef/code-gen/package-update"
	"log/slog"
	"os"
	"strings"
)

const corePathPackage = "github.com/ihatiko/chef/cli/go-chef-core"
const coreNamePackage = "go-chef-core"

func main() {
	params := strings.Join(os.Args[1:], " ")
	//TODO timeout on update
	packageUpdate.AutoUpdate(corePathPackage)
	composer := filebuilder.NewComposer()
	result, err := composer.ExecDefaultCommand(coreNamePackage)
	if err != nil {
		slog.Error("Error executing command: ", slog.Any("error", err), slog.String("command", params))
	}
	slog.Info(result.String())
}
