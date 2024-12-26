package main

import (
	"fmt"
	filebuilder "github.com/ihatiko/chef/code-gen-utils/file-builder"
	packageUpdate "github.com/ihatiko/chef/code-gen-utils/package-update"
	"log/slog"
	"os"
	"strings"
)

const corePathPackage = "github.com/ihatiko/chef/cli/go-chef-core"
const coreNamePackage = "go-chef-core"

func main() {
	fmt.Println(os.Args)
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
