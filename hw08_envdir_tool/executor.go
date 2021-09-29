package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) // nolint: gosec
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := applyEnvs(env)
	if err != nil {
		fmt.Printf("apply env config: %v", err)
		return 1
	}
	command.Env = append(command.Env, os.Environ()...)

	if err := command.Start(); err != nil {
		log.Fatalf("cmd start error: %v", err)
	}

	if err := command.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok { // nolint: errorlint
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		} else {
			fmt.Printf("command wait: %v", err)
			return 1
		}
	}

	return 0
}

func applyEnvs(env Environment) error {
	for k, v := range env {
		err := os.Unsetenv(k)
		if err != nil {
			return err
		}
		if v.NeedRemove {
			continue
		}
		err = os.Setenv(k, v.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
