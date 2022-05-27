package main

import (
	"errors"
	"os"
	"os/exec"
)

const (
	exitCodeOk             = 0
	exitCodeErr            = 1
	exitCodeCannotUnsetEnv = 100
	exitCodeCannotSetEnv   = 101
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for name, ev := range env {
		if ev.NeedRemove {
			if err := os.Unsetenv(name); err != nil {
				return exitCodeCannotUnsetEnv
			}

			continue
		}

		if err := os.Setenv(name, ev.Value); err != nil {
			return exitCodeCannotSetEnv
		}
	}

	run, args := cmd[0], cmd[1:]

	exe := exec.Command(run, args...)

	exe.Stdin = os.Stdin
	exe.Stdout = os.Stdout
	exe.Stderr = os.Stderr

	if err := exe.Run(); err != nil {
		var exitError *exec.ExitError

		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		return exitCodeErr
	}

	return exitCodeOk
}
