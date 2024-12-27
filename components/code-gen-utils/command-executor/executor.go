package command_executor

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

func NewCommand(cmd string, skipOnError bool) Command {
	return Command{cmd: cmd, skipOnError: skipOnError, state: true}
}

func NewDefaultCommand(cmd string) Command {
	return Command{cmd: cmd, skipOnError: true, state: true}
}

func NewConditionalCommand(cmd string, state bool, skipOnError bool) Command {
	return Command{cmd: cmd, skipOnError: skipOnError, state: state}
}

type Command struct {
	cmd         string
	skipOnError bool
	state       bool
}

type Executor struct {
	ProjectPath string
}

func NewExecutor(projectPaths ...string) *Executor {
	projectPath := ""
	if projectPaths != nil || len(projectPaths) > 0 {
		projectPath = projectPaths[0]
	}
	return &Executor{ProjectPath: GetPath(projectPath)}
}

func (c *Executor) CommandsExecutor(commands ...Command) {
	consoleEnv := "bash"
	if os.Getenv("GOOS") == "windows" || strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		consoleEnv = "powershell"
	}
	for _, command := range commands {
		if command.state {
			_, err := c.ExecCommand(command.cmd, consoleEnv)
			slog.Info(command.cmd)
			if err != nil && !command.skipOnError {
				slog.Error(err.Error())
				break
			}
		}
	}
}

func (c *Executor) ConditionalExecutor(commands ...Command) {
	consoleEnv := "bash"
	if os.Getenv("GOOS") == "windows" || strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		consoleEnv = "powershell"
	}
	for _, command := range commands {
		_, err := c.ExecCommand(command.cmd, consoleEnv)
		slog.Info(command.cmd)
		if err != nil && !command.skipOnError {
			slog.Error(err.Error())
			break
		}
	}
}
func (c *Executor) ExecDefaultCommand(command string) (*strings.Builder, error) {
	consoleEnv := "bash"
	if os.Getenv("GOOS") == "windows" || strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		consoleEnv = "powershell"
	}
	return c.ExecCommand(command, consoleEnv)
}

func (c *Executor) ExecCommand(command string, consoleEnv string) (*strings.Builder, error) {
	cmdFolder := exec.Command(consoleEnv, "-c", command)
	builder := new(strings.Builder)
	builderErr := new(strings.Builder)
	cmdFolder.Stdin = strings.NewReader("")
	cmdFolder.Stderr = builderErr
	cmdFolder.Dir = c.ProjectPath
	cmdFolder.Stdout = builder
	err := cmdFolder.Run()
	if err != nil {
		return builder, err
	}
	if cmdFolder.ProcessState.ExitCode() > 0 {
		return builder, errors.New(builderErr.String())
	}
	return builder, err
}
