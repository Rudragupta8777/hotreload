package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// 1. Define the command-line flags
	rootFlag := flag.String("root", ".", "Directory to watch for file changes")
	buildFlag := flag.String("build", "", "Command used to build the project")
	execFlag := flag.String("exec", "", "Command used to run the built server")

	// Parse the flags from the terminal
	flag.Parse()

	// 2. Set up slog for structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// 3. Validate the inputs
	if *buildFlag == "" || *execFlag == "" {
		slog.Error("Both --build and --exec commands are required.")
		os.Exit(1)
	}

	slog.Info("Starting hotreload engine",
		"root", *rootFlag,
		"build_cmd", *buildFlag,
		"exec_cmd", *execFlag,
	)

	// Initialize the File Watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("Failed to create watcher", "error", err)
		os.Exit(1)
	}
	defer watcher.Close()

	// Add all subdirectories to the watcher
	watchDirectories(watcher, *rootFlag)

	// Initialize the Runner
	runner := NewRunner(*buildFlag, *execFlag)

	// Trigger the very first build immediately (Requirement)
	go runner.TriggerBuildAndRun()

	// Debounce logic parameters
	debounceDuration := 500 * time.Millisecond
	var timer *time.Timer

	slog.Info("Listening for changes...")

	// Infinite loop to listen for file events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// We only care about Writes, Creates, or Removes
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {

				// Bonus point: If a developer creates a new folder, watch it dynamically!
				if event.Has(fsnotify.Create) {
					info, err := os.Stat(event.Name)
					if err == nil && info.IsDir() && !isIgnored(event.Name) {
						watcher.Add(event.Name)
						slog.Info("Detected and watching new directory", "path", event.Name)
					}
				}

				// Debounce: Reset the timer every time a new event comes in quickly
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(debounceDuration, func() {
					slog.Info("Changes detected, reloading...", "file", event.Name)
					runner.TriggerBuildAndRun()
				})
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			slog.Error("Watcher error", "error", err)
		}
	}
}