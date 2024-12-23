package utils

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

func NewComposer(projectPath string) *Composer {
	return &Composer{ProjectPath: GetPath(projectPath)}
}

func (c *Composer) CommandsComposer(commands ...Command) {
	consoleEnv := "bash"
	if os.Getenv("GOOS") == "windows" || strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		consoleEnv = "powershell"
	}
	for _, command := range commands {
		if command.state {
			err := c.ExecCommand(command.cmd, consoleEnv)
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
		err := c.ExecCommand(command.cmd, consoleEnv)
		slog.Info(command.cmd)
		if err != nil && !command.skipOnError {
			slog.Error(err.Error())
			break
		}
	}
}

func (c *Composer) ExecCommand(command string, consoleEnv string) error {
	cmdFolder := exec.Command(consoleEnv, "-c", command)
	var out bytes.Buffer
	cmdFolder.Stdin = strings.NewReader("some input")
	cmdFolder.Stderr = &out
	cmdFolder.Dir = c.ProjectPath
	err := cmdFolder.Run()
	if err != nil {
		return err
	}
	if cmdFolder.ProcessState.ExitCode() > 0 {
		return errors.New(out.String())
	}
	return nil
}
