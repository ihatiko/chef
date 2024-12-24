package file_builder

import (
	"bytes"
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

type Composer struct {
	ProjectPath string
}

func NewComposer(projectPaths ...string) *Composer {
	projectPath := ""
	if projectPaths != nil || len(projectPaths) > 0 {
		projectPath = projectPaths[0]
	}
	return &Composer{ProjectPath: GetPath(projectPath)}
}

func (c *Composer) CommandsComposer(commands ...Command) {
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

func (c *Composer) ConditionalComposer(commands ...Command) {
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
func (c *Composer) ExecDefaultCommand(command string) (*strings.Builder, error) {
	consoleEnv := "bash"
	if os.Getenv("GOOS") == "windows" || strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		consoleEnv = "powershell"
	}
	return c.ExecCommand(command, consoleEnv)
}

func (c *Composer) ExecCommand(command string, consoleEnv string) (*strings.Builder, error) {
	cmdFolder := exec.Command(consoleEnv, "-c", command)
	builder := new(strings.Builder)
	var out bytes.Buffer
	cmdFolder.Stdin = strings.NewReader("")
	cmdFolder.Stderr = &out
	cmdFolder.Dir = c.ProjectPath
	cmdFolder.Stdout = builder
	err := cmdFolder.Run()
	if err != nil {
		return builder, err
	}
	if cmdFolder.ProcessState.ExitCode() > 0 {
		return builder, errors.New(out.String())
	}
	return builder, err
}
