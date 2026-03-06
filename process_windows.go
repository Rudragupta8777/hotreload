//go:build windows

package main

import (
	"fmt"
	"os/exec"
)

// killProcessTree forcefuly kills the process and all its children on Windows
func killProcessTree(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		killCmd := exec.Command("taskkill", "/T", "/F", "/PID", fmt.Sprint(cmd.Process.Pid))
		killCmd.Run()
	}
}

// setupProcessGroup is a no-op on Windows because taskkill handles the tree
func setupProcessGroup(cmd *exec.Cmd) {
	// Do nothing
}