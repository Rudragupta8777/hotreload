//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

// killProcessTree kills the entire process group on Unix/Linux/Mac
func killProcessTree(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
}

// setupProcessGroup assigns the process to a new group so we can kill its children later
func setupProcessGroup(cmd *exec.Cmd) {
	if cmd != nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}
}