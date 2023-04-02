package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envKey := range env {
		if env[envKey].NeedRemove {
			if _, exists := os.LookupEnv(envKey); exists {
				os.Unsetenv(envKey)
			}
		} else {
			os.Setenv(envKey, env[envKey].Value)
		}
	}

	cmdX := exec.Command(cmd[0], cmd[1:]...)
	cmdX.Env = os.Environ()
	cmdX.Stderr = os.Stderr
	cmdX.Stdin = os.Stdin
	cmdX.Stdout = os.Stdout

	err := cmdX.Run()
	if err != nil {
		err2, ok := err.(*exec.ExitError)
		if ok {
			return err2.ProcessState.ExitCode()
		} else {
			return -1
		}
	}
	return 0
}
