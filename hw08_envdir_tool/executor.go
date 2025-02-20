package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	err := fillEnv(env)
	if err != nil {
		slog.Error(fmt.Errorf("cmd run err: %w", err).Error())
		return
	}

	path, err := exec.LookPath(cmd[0])
	if err != nil {
		slog.Error(fmt.Errorf("cmd run err: %w", err).Error())
	}

	command := exec.Command(path, cmd[1:]...)

	command.Stdin = os.Stdout
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		slog.Error(fmt.Errorf("cmd run err: %w", err).Error())
		return command.ProcessState.ExitCode()
	}

	return command.ProcessState.ExitCode()
}

func fillEnv(env Environment) error {
	for _, val := range os.Environ() {
		parts := strings.SplitN(val, "=", 2)
		if len(parts) == 2 {
			err := os.Setenv(parts[0], parts[1])
			if err != nil {
				return fmt.Errorf("set env err: %w", err)
			}
		}
	}

	for key, val := range env {
		if val.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return fmt.Errorf("unset env err: %w", err)
			}

			continue
		}

		_, ok := os.LookupEnv(key)
		if ok {
			err := os.Unsetenv(key)
			if err != nil {
				return fmt.Errorf("unset env err: %w", err)
			}
		}

		err := os.Setenv(key, val.Value)
		if err != nil {
			return fmt.Errorf("set env err: %w", err)
		}
	}

	return nil
}
