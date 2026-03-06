package main

import (
	"flag"
	"log/slog"
	"os"
)

func main() {
	// 1. Defined the command-line flags
	rootFlag := flag.String("root", ".", "Directory to watch for file changes")
	buildFlag := flag.String("build", "", "Command used to build the project")
	execFlag := flag.String("exec", "", "Command used to run the built server")

	// Parse the flags from the terminal
	flag.Parse()

	// 2. Set up slog for structured logging
	// We use a TextHandler here which is great for CLI tools.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // We'll set this to debug while building
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

	
}