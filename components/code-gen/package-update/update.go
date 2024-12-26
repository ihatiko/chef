package package_update

import (
	"bufio"
	"fmt"
	filebuilder "github.com/ihatiko/chef/code-gen-utils/file-builder"
	"golang.org/x/mod/semver"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func AutoUpdate(packageName string) {
	splittedPackage := strings.Split(packageName, "/")
	packageUrl := splittedPackage[len(splittedPackage)-1]
	composer := filebuilder.NewComposer()
	currentVersion, err := composer.ExecDefaultCommand(fmt.Sprintf("%s version", packageUrl))
	if err != nil {
		slog.Error("Failed to execute composer", slog.Any("error", err), slog.String("package", packageUrl))
		return
	}
	fullPathName := fmt.Sprintf("https://proxy.golang.org/%s/@v/list", packageName)
	response, err := http.Get(fullPathName)
	if err != nil {
		slog.Error("Error fetching latest version of package", slog.Any("error", err))
		return
	}

	reader := bufio.NewReader(response.Body)

	bytes, err := reader.ReadBytes(0)
	if err != nil && err != io.EOF {
		slog.Error("Error reading response", slog.Any("error", err))
		return
	}

	versions := strings.Split(string(bytes), "\n")

	semver.Sort(versions)
	lastVersion := versions[len(versions)-1]
	formattedCurrentVersion := strings.ReplaceAll(currentVersion.String(), "\n", "")
	if lastVersion == formattedCurrentVersion {
		slog.Info("actual", slog.String("package", packageName), slog.String("version", formattedCurrentVersion))
		return
	}
	slog.Info("try update", slog.String("current-version", formattedCurrentVersion), slog.String("last-version", lastVersion))
	installInstruction := fmt.Sprintf("%s@%s", packageName, lastVersion)
	command := fmt.Sprintf("go install %s", installInstruction)

	slog.Info("Executing", slog.String("command", command))
	_, err = composer.ExecDefaultCommand(command)
	if err != nil {
		slog.Error("error executing composer", slog.Any("error", err))
		return
	}
}
