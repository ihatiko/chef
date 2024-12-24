package main

import (
	"bufio"
	"fmt"
	filebuilder "github.com/ihatiko/chef/code-gen-utils/file-builder"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// TODO сделать подгрузку динамических модулей через базовый конфиг
// Использовать данный слой только как фронтенд и хранение внешних настроек
var rootCmd = &cobra.Command{
	Use:   "zero",
	Short: "zero is a cli tool for performing basic mathematical operations",
	Long:  "zero is a cli tool for performing basic mathematical operations - addition, multiplication, division and subtraction.",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Zero '%s'\n", err)
		os.Exit(1)
	}
}
func AutoUpdate(packageName string) {
	pkgs := strings.Split(packageName, "/")
	packageUrl := pkgs[len(pkgs)-1]
	composer := filebuilder.NewComposer()
	currentVersion, err := composer.ExecDefaultCommand(fmt.Sprintf("%s version", packageUrl))
	if err != nil {
		slog.Error("Failed to execute composer", slog.Any("error", err))
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
func main() {
	AutoUpdate("github.com/ihatiko/go-chef-sandbox-test")
}
