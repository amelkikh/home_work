package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:", os.Args[0], "/path/to/env/dir", "command", "arg1", "argN")
		return
	}
	envDir := os.Args[1]
	envs, err := ReadDir(envDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Exit(RunCmd(os.Args[2:], envs))
}
