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
	mu       sync.Mutex // Prevents race conditions if builds trigger rapidly
}

func NewRunner(buildCmd, execCmd string) *Runner {
	return &Runner{
		buildCmd: buildCmd,
		execCmd:  execCmd,
	}
}

// Kill stops the currently running server
func (r *Runner) Kill() {
	if r.cancel != nil {
		r.cancel() // Canceling the context signals the process to stop
	}
	if r.cmd != nil && r.cmd.Process != nil {
		slog.Info("Terminating previous process...")
		r.cmd.Process.Kill()
		r.cmd.Wait() // Wait for it to fully die so we free up the port
	}
}

// TriggerBuildAndRun kills the old process, builds the new code, and runs it
func (r *Runner) TriggerBuildAndRun() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Kill()

	// Create a new context for the upcoming processes
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel

	// 1. Build Phase
	slog.Info("Building project...", "command", r.buildCmd)
	buildArgs := strings.Fields(r.buildCmd)
	build := exec.CommandContext(ctx, buildArgs[0], buildArgs[1:]...)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	if err := build.Run(); err != nil {
		slog.Error("Build failed! Waiting for next file change...", "error", err)
		return // Stop here. Don't try to run a broken build.
	}
	slog.Info("Build successful.")

	// 2. Run Phase
	slog.Info("Starting server...", "command", r.execCmd)
	execArgs := strings.Fields(r.execCmd)
	r.cmd = exec.CommandContext(ctx, execArgs[0], execArgs[1:]...)
	
	// Stream logs in real-time (Requirement)
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr

	if err := r.cmd.Start(); err != nil {
		slog.Error("Failed to start server", "error", err)
		return
	}
}