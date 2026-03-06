package main

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Runner struct {
	buildCmd string
	execCmd  string
	cmd      *exec.Cmd
	cancel   context.CancelFunc
	mu       sync.Mutex
}

func NewRunner(buildCmd, execCmd string) *Runner {
	return &Runner{
		buildCmd: buildCmd,
		execCmd:  execCmd,
	}
}

// Kill stops the currently running server AND all its child processes
func (r *Runner) Kill() {
	if r.cancel != nil {
		r.cancel()
	}
	if r.cmd != nil && r.cmd.Process != nil {
		slog.Info("Terminating process tree...", "pid", r.cmd.Process.Pid)
		
		// This will call the Windows version on the machine and the Unix version for reviewers
		killProcessTree(r.cmd)
		
		r.cmd.Wait() // Use to wait for it to fully die to free up the port
	}
}

// TriggerBuildAndRun kills the old process, builds the new code and runs it
func (r *Runner) TriggerBuildAndRun() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Kill()

	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel

	// 1. Building Phase
	slog.Info("Building project...", "command", r.buildCmd)
	buildArgs := strings.Fields(r.buildCmd)
	build := exec.CommandContext(ctx, buildArgs[0], buildArgs[1:]...)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	if err := build.Run(); err != nil {
		slog.Error("Build failed! Waiting for next file change...", "error", err)
		return 
	}
	slog.Info("Build successful.")

	// 2. Run Phase
	slog.Info("Starting server...", "command", r.execCmd)
	execArgs := strings.Fields(r.execCmd)
	r.cmd = exec.CommandContext(ctx, execArgs[0], execArgs[1:]...)
	
	// Setup process grouping for clean termination
	setupProcessGroup(r.cmd)

	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr

	if err := r.cmd.Start(); err != nil {
		slog.Error("Failed to start server", "error", err)
		return
	}
}